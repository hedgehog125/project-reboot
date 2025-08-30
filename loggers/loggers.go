package loggers

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"reflect"
	"slices"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

const (
	PublicMessageKey = "publicMessage"
	UserIDKey        = "userID"

	MaxSaveBatchSize = 100
)

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
	shutdownFinishedChan chan struct{}
}
type entry struct {
	time           time.Time
	timeKnown      bool
	level          int
	message        string
	attributes     map[string]any
	sourceFile     string
	sourceFunction string
	sourceLine     int
	publicMessage  string
	userID         int
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
			Level: level,
		}),
		baseAttrs:            map[string]any{},
		entryChan:            make(chan *entry, 100),
		requestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
	}
}

func (handler *Handler) Listen() {
	shuttingDown := false
listenLoop:
	for {
		entries := []*entry{}
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

		select {
		case entry := <-handler.entryChan:
			entries = append(entries, entry)
		case <-handler.requestShutdownChan:
			shuttingDown = true
			emptyEntryChan()
			break listenLoop
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
					break collectBatchLoop
				}
				if len(entries) >= MaxSaveBatchSize {
					break
				}
			}
		}

		stdErr := dbcommon.WithWriteTx(
			context.TODO(), handler.App.Database,
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
				}).Exec(ctx)
			},
		)
		if stdErr != nil {
			// TODO: fall back to individual creates, some entries may be invalid
			// TODO: log a warning if they all succeed, because otherwise doing this just reduced performance for no benefit

			// TODO: can the logger call itself?
			fmt.Printf(
				"error: unable to store logs to database. error:\n%v",
				stdErr.Error(),
			)
		}

		if shuttingDown {
			break listenLoop
		}
	}
	close(handler.requestShutdownChan)
	close(handler.shutdownFinishedChan)
}
func (handler *Handler) Shutdown() {
	// TODO: timeout?
	// TODO: what if it's not running?
	fmt.Println("logger shutting down")
	// TODO: this panics if the handler has to shut itself down due to an error
	handler.requestShutdownChan <- struct{}{}
	fmt.Println("logger finishing saving logs")
	<-handler.shutdownFinishedChan
	fmt.Println("logger stopped")
}

func (handler Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= handler.Level
}

func (handler Handler) Handle(ctx context.Context, record slog.Record) error {
	entry := &entry{
		level:   int(record.Level),
		message: record.Message,
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
	resolvedAttrs, specialProps := handler.resolveNestedAttrs(attrs)
	entry.publicMessage = specialProps.publicMessage
	entry.userID = specialProps.userID
	entry.attributes = resolvedAttrs

	stdErr := handler.textHandler.Handle(ctx, record)
	if stdErr != nil {
		fmt.Printf(
			"warning: log Handler.textHandler.Handle returned an error. error:\n%v",
			stdErr.Error(),
		)
	}
	handler.entryChan <- entry
	return nil
}

type specialProperties struct {
	publicMessage string
	userID        int
}

func (handler Handler) resolveNestedAttrs(attrs []slog.Attr) (map[string]any, specialProperties) {
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
		newSpecialProps := handler.appendAttr(attr, nestedResolved, isTopLevel)
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
func (handler Handler) appendAttr(attr slog.Attr, outputAttrs map[string]any, isTopLevel bool) specialProperties {
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
				newSpecialProps := handler.appendAttr(attr, outputAttrs, true)
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
				_ = handler.appendAttr(attr, groupAttr, false)
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
			intValue, ok := attr.Value.Any().(int)
			if ok {
				specialProps.userID = intValue
			} else {
				typeOf := reflect.TypeOf(attr.Value.Any())
				fmt.Printf("warning: userID property in log statement was not an int so has been ignored. type: %v", typeOf) // TODO: can the logger call itself?
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
	resolvedAttrs, specialProps := handler.resolveNestedAttrs(attrs)
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
