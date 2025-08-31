package loggers

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"os"
	"reflect"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers"
)

const (
	// Special attributes
	PublicMessageKey = "publicMessage"
	UserIDKey        = "userID"

	MaxSaveBatchSize = 100
	ShutdownTimeout  = 5 * time.Second
)

type disableErrorLoggingKey = struct{} // Used to prevent infinite loops

type Handler struct {
	App                  *common.App
	Level                slog.Level
	SaveToDatabase       bool
	ShouldPrint          bool
	textHandler          slog.Handler
	baseAttrs            map[string]any
	baseSpecialProps     specialProperties
	baseGroups           []string
	entryChan            chan *entry
	requestShutdownChan  chan struct{}
	shutdownCtx          context.Context
	cancelShutdownCtx    context.CancelFunc
	shutdownFinishedChan chan struct{}
	mu                   *sync.Mutex
}
type entry struct {
	time                time.Time
	timeKnown           bool
	level               int
	message             string
	attributes          map[string]any
	sourceFile          string
	sourceFunction      string
	sourceLine          int
	publicMessage       string
	userID              int
	disableErrorLogging bool
}

func NewHandler(
	level slog.Level, saveToDatabase bool, shouldPrint bool,
	app *common.App,
) Handler {
	return Handler{
		App:            app,
		Level:          level,
		SaveToDatabase: saveToDatabase,
		ShouldPrint:    shouldPrint,
		textHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: true,
		}),
		baseAttrs:            map[string]any{},
		entryChan:            make(chan *entry, 100),
		requestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
		mu:                   &sync.Mutex{},
	}
}

func (handler *Handler) Listen() {
	shuttingDown := false
	loggedBulkWarning := false
	loggedAdminNotificationError := false
listenLoop:
	for {
		shouldReEnableSelfLogging := false
		entries := []*entry{}
		selfLogged := false
		emptyEntryChan := func() {
			for {
				select {
				case entry := <-handler.entryChan:
					entries = append(entries, entry)
				default:
					return
				}
			}
		}

		if shuttingDown {
			emptyEntryChan()
		} else {
			select {
			case entry := <-handler.entryChan:
				entries = append(entries, entry)
			case <-handler.requestShutdownChan:
				shuttingDown = true
				emptyEntryChan()
				break listenLoop
			}
		}

		if !shuttingDown {
			timeoutChan := time.After(handler.App.Env.LOG_STORE_INTERVAL)
		collectBatchLoop:
			for {
				select {
				case entry := <-handler.entryChan:
					entries = append(entries, entry)
				case <-handler.requestShutdownChan:
					shuttingDown = true
					emptyEntryChan()
					break collectBatchLoop
				case <-timeoutChan:
					shouldReEnableSelfLogging = true
					break collectBatchLoop
				}
				if len(entries) >= MaxSaveBatchSize {
					break
				}
			}
		}

		bulkWriteErr := handler.bulkWrite(entries)
		if bulkWriteErr != nil {
			if handler.individualWriteFallback(entries, bulkWriteErr, &loggedBulkWarning) {
				shouldReEnableSelfLogging = false
				selfLogged = true
			}
		}
		if handler.maybeNotifyAdmin(entries, &loggedAdminNotificationError) {
			shouldReEnableSelfLogging = false
			selfLogged = true
		}

		if shouldReEnableSelfLogging {
			loggedBulkWarning = false
			loggedAdminNotificationError = false
		}
		if shuttingDown {
			if selfLogged {
				select {
				case <-handler.shutdownCtx.Done():
					break listenLoop
				default:
				}
				time.Sleep(5 * time.Millisecond) // Give the channels a second so that all the entries that were added can be read
			} else {
				break listenLoop
			}
		}
	}
	close(handler.requestShutdownChan)
	close(handler.shutdownFinishedChan)
	handler.cancelShutdownCtx()
}
func (handler *Handler) bulkWrite(entries []*entry) error {
	ctx := handler.shutdownCtx
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
	}
	return dbcommon.WithWriteTx(
		ctx, handler.App.Database,
		func(tx *ent.Tx, ctx context.Context) error {
			return tx.LogEntry.MapCreateBulk(entries, func(lec *ent.LogEntryCreate, i int) {
				entry := entries[i]
				lec.SetTime(entry.time).SetTimeKnown(entry.timeKnown).
					SetLevel(entry.level).
					SetMessage(entry.message).
					SetAttributes(entry.attributes).
					SetSourceFile(entry.sourceFile).
					SetSourceFunction(entry.sourceFunction).
					SetSourceLine(entry.sourceLine).
					SetPublicMessage(entry.publicMessage)
				if entry.userID != 0 {
					lec.SetUserID(entry.userID)
				}
			}).Exec(ctx)
		},
	)
}
func (handler *Handler) individualWriteFallback(
	entries []*entry,
	bulkWriteErr error,
	loggedBulkWarningPtr *bool,
) bool {
	selfLogged := false
	allSucceeded := true
	for _, entry := range entries {
		var timeout time.Duration
		if entry.level >= int(slog.LevelError) {
			timeout = time.Second
		} else if entry.level >= int(slog.LevelWarn) {
			timeout = 500 * time.Millisecond
		} else {
			timeout = 100 * time.Millisecond
		}
		baseCtx := context.Background()
		if handler.shutdownCtx != nil {
			baseCtx = handler.shutdownCtx
		}
		ctx, cancel := context.WithTimeout(baseCtx, timeout)
		stdErr := dbcommon.WithWriteTx(
			ctx, handler.App.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				builder := tx.LogEntry.Create().
					SetTime(entry.time).SetTimeKnown(entry.timeKnown).
					SetLevel(entry.level).
					SetMessage(entry.message).
					SetAttributes(entry.attributes).
					SetSourceFile(entry.sourceFile).
					SetSourceFunction(entry.sourceFunction).
					SetSourceLine(entry.sourceLine).
					SetPublicMessage(entry.publicMessage)
				if entry.userID != 0 {
					// TODO: don't set this and instead hydrate in a second query, that way invalid IDs don't prevent log storage
					builder.SetUserID(entry.userID)
				}
				return builder.Exec(ctx)
			},
		)
		cancel()
		if stdErr != nil {
			allSucceeded = false
			if !entry.disableErrorLogging {
				pc, _, _, _ := runtime.Caller(0)
				record := slog.NewRecord(
					handler.App.Clock.Now(),
					slog.LevelError,
					"failed to write log entry to database",
					pc,
				)
				record.AddAttrs(slog.Any("log", entry))
				record.AddAttrs(slog.Any("error", stdErr))
				handler.Handle(
					context.WithValue(context.Background(), disableErrorLoggingKey{}, true),
					record,
				)
				selfLogged = true
			}
		}
	}
	if allSucceeded && !*loggedBulkWarningPtr {
		pc, _, _, _ := runtime.Caller(0)
		record := slog.NewRecord(
			handler.App.Clock.Now(),
			slog.LevelWarn,
			"bulk log write failed but the individual fallback writes all succeeded, so the writes took longer than they should have",
			pc,
		)
		record.AddAttrs(slog.Any("error", bulkWriteErr))
		handler.Handle(
			context.Background(),
			record,
		)
		*loggedBulkWarningPtr = true
		selfLogged = true
	}
	return selfLogged
}
func (handler *Handler) maybeNotifyAdmin(entries []*entry, loggedAdminNotificationErrorPtr *bool) bool {
	if *loggedAdminNotificationErrorPtr {
		return false
	}
	selfLogged := false

	shouldNotifyAdmin := false
	if handler.App.Env.ADMIN_USERNAME != "" {
		for _, entry := range entries {
			if entry.level >= int(slog.LevelError) {
				shouldNotifyAdmin = true
				break
			}
		}
	}
	if shouldNotifyAdmin {
		// TODO: rate limiting!

		// TODO: reserve a bit of time for this in case the database writing times out during a shutdown
		baseCtx := context.Background()
		if handler.shutdownCtx != nil {
			baseCtx = handler.shutdownCtx
		}
		ctx, cancel := context.WithTimeout(baseCtx, 2*time.Second)
		var queuedCount int
		var errs map[string]*common.Error
		stdErr := dbcommon.WithWriteTx(
			ctx, handler.App.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				userOb, stdErr := tx.User.Query().Where(user.Username(handler.App.Env.ADMIN_USERNAME)).Only(ctx)
				if stdErr != nil {
					return stdErr
				}
				var commErr *common.Error
				queuedCount, errs, commErr = handler.App.Messengers.SendUsingAll(
					&common.Message{
						Type: common.MessageAdminError,
						User: userOb,
					},
					ctx,
				)

				return commErr.StandardError()
			},
		)
		cancel()
		if stdErr != nil {
			// TODO: add a special context item to notify the admin by crashing
			pc, _, _, _ := runtime.Caller(0)
			record := slog.NewRecord(
				handler.App.Clock.Now(),
				slog.LevelError,
				"failed to message admin about an error",
				pc,
			)
			record.AddAttrs(slog.Any("error", stdErr))
			handler.Handle(
				context.Background(),
				record,
			)
			*loggedAdminNotificationErrorPtr = true
			return true
		}

		if len(errs) > 0 { // SendUsingAll will have logged
			// TODO: log something with a special context item to notify the admin by crashing
			*loggedAdminNotificationErrorPtr = true
			selfLogged = true
		}
		if queuedCount == 0 {
			// TODO: add a special context item to notify the admin by crashing
			message := "admin user has no contacts so couldn't notify them about an error"
			for _, commErr := range errs {
				if !errors.Is(commErr, messengers.ErrNoContactForUser) {
					message = "unable to prepare messages to notify admin about an error, see the errors before"
				}
			}

			pc, _, _, _ := runtime.Caller(0)
			record := slog.NewRecord(
				handler.App.Clock.Now(),
				slog.LevelError,
				message,
				pc,
			)
			handler.Handle(
				context.Background(),
				record,
			)
			*loggedAdminNotificationErrorPtr = true
			selfLogged = true
		}
	}
	return selfLogged
}

func (handler *Handler) Shutdown() {
	// TODO: what if it's not running?
	handler.mu.Lock()
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	handler.shutdownCtx = ctx
	handler.cancelShutdownCtx = cancel
	handler.mu.Unlock()
	handler.requestShutdownChan <- struct{}{}
	<-handler.shutdownFinishedChan
}

func (handler Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= handler.Level
}

func (handler Handler) Handle(ctx context.Context, record slog.Record) error {
	disableErrLogging, _ := ctx.Value(disableErrorLoggingKey{}).(bool)
	entry := &entry{
		level:               int(record.Level),
		message:             record.Message,
		disableErrorLogging: disableErrLogging,
	}
	if !record.Time.IsZero() {
		entry.time = record.Time
		entry.timeKnown = true
	}
	source := record.Source()
	if source != nil {
		entry.sourceFile = source.File
		entry.sourceFunction = source.Function
		entry.sourceLine = source.Line
	}

	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	resolvedAttrs, specialProps := handler.resolveNestedAttrs(attrs, !disableErrLogging)
	entry.publicMessage = specialProps.publicMessage
	entry.userID = specialProps.userID
	entry.attributes = resolvedAttrs

	stdErr := handler.textHandler.Handle(ctx, record)
	handler.entryChan <- entry
	if stdErr != nil && !disableErrLogging {
		pc, _, _, _ := runtime.Caller(0)
		record := slog.NewRecord(
			handler.App.Clock.Now(),
			slog.LevelWarn,
			"logger Handler.textHandler.Handle returned an error",
			pc,
		)
		record.AddAttrs(slog.Any("error", stdErr))
		handler.Handle(
			context.WithValue(context.Background(), disableErrorLoggingKey{}, true),
			record,
		)
	}
	return nil
}

type specialProperties struct {
	publicMessage string
	userID        int
}

func (handler Handler) resolveNestedAttrs(attrs []slog.Attr, logErrors bool) (map[string]any, specialProperties) {
	resolved := maps.Clone(handler.baseAttrs)
	nestedResolved := resolved
	for _, key := range handler.baseGroups {
		newMap, ok := nestedResolved[key].(map[string]any)
		if ok {
			newMap = maps.Clone(newMap)
		} else {
			newMap = map[string]any{}
		}
		nestedResolved[key] = newMap
		nestedResolved = newMap
	}

	isTopLevel := len(handler.baseGroups) == 0
	specialProps := specialProperties{
		publicMessage: handler.baseSpecialProps.publicMessage,
		userID:        handler.baseSpecialProps.userID,
	}
	for _, attr := range attrs {
		newSpecialProps := handler.appendAttr(attr, nestedResolved, isTopLevel, logErrors)
		if newSpecialProps.publicMessage != "" {
			specialProps.publicMessage = newSpecialProps.publicMessage
		}
		if newSpecialProps.userID != 0 {
			specialProps.userID = newSpecialProps.userID
		}
	}
	return resolved, specialProps
}

// Note: handler.baseGroups is handled by appendNestedAttrs instead
func (handler Handler) appendAttr(
	attr slog.Attr, outputAttrs map[string]any,
	isTopLevel bool, logErrors bool,
) specialProperties {
	specialProps := specialProperties{}
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return specialProps
	}

	kind := attr.Value.Kind()
	if kind == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		if len(groupAttrs) == 0 {
			return specialProps
		}
		// If the key is non-empty, write it out and indent the rest of the attrs.
		// Otherwise, inline the attrs.
		if attr.Key == "" { // Inline
			for _, attr := range groupAttrs {
				newSpecialProps := handler.appendAttr(attr, outputAttrs, true, logErrors)
				if newSpecialProps.publicMessage != "" {
					specialProps.publicMessage = newSpecialProps.publicMessage
				}
				if newSpecialProps.userID != 0 {
					specialProps.userID = newSpecialProps.userID
				}
			}
		} else {
			groupAttr := map[string]any{}
			for _, attr := range groupAttrs {
				_ = handler.appendAttr(attr, groupAttr, false, logErrors)
			}
			outputAttrs[attr.Key] = groupAttr
		}
		return specialProps
	}
	if isTopLevel {
		if attr.Key == PublicMessageKey {
			specialProps.publicMessage = attr.Value.String()
			return specialProps
		}
		if attr.Key == UserIDKey {
			intValue, ok := attr.Value.Any().(int64)
			if ok {
				specialProps.userID = int(intValue)
			} else {
				if logErrors {
					pc, _, _, _ := runtime.Caller(0)
					record := slog.NewRecord(
						handler.App.Clock.Now(),
						slog.LevelWarn,
						"userID property in log statement was not an int so has been ignored",
						pc,
					)
					record.AddAttrs(slog.String("type", reflect.TypeOf(attr.Value.Any()).String()))
					handler.Handle(
						context.WithValue(context.Background(), disableErrorLoggingKey{}, true),
						record,
					)
				}
			}
			outputAttrs[attr.Key] = attr.Value.Any() // Also store the value in the attributes so it's preserved if the user is deleted
			return specialProps
		}
	}

	outputAttrs[attr.Key] = attr.Value.Any()
	return specialProps
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (handler Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return handler
	}
	handler.textHandler = handler.textHandler.WithAttrs(attrs)
	resolvedAttrs, specialProps := handler.resolveNestedAttrs(attrs, true)
	// maps.Copy(handler.baseAttrs, resolvedAttrs) // Mutate baseAttrs rather than copying so other references are updated
	handler.baseAttrs = resolvedAttrs

	if specialProps.publicMessage != "" {
		handler.baseSpecialProps.publicMessage = specialProps.publicMessage
	}
	if specialProps.userID != 0 {
		handler.baseSpecialProps.userID = specialProps.userID
	}

	return handler
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (handler Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return handler
	}
	handler.textHandler = handler.textHandler.WithGroup(name)
	handler.baseGroups = slices.Concat(handler.baseGroups, []string{name})
	// handler.baseAttrs = map[string]any{}
	return handler
}
