package loggers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
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
		entryChan:            make(chan *entry, 100),
		requestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
	}
}

func (handler *Handler) Listen() {
listenLoop:
	for {
		entries := []*entry{}
		select {
		case entry := <-handler.entryChan:
			entries = append(entries, entry)
		case <-handler.requestShutdownChan:
			break listenLoop
		}

		timeoutChan := time.After(handler.App.Env.LOG_STORE_INTERVAL)
	collectBatchLoop:
		for {
			select {
			case entry := <-handler.entryChan:
				entries = append(entries, entry)
			case <-handler.requestShutdownChan:
				for {
					select {
					case entry := <-handler.entryChan:
						entries = append(entries, entry)
					default:
						break collectBatchLoop
					}
				}
			case <-timeoutChan:
				break collectBatchLoop
			}
			if len(entries) >= MaxSaveBatchSize {
				break
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

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
//
// [Logger] discards any errors from Handle. Wrap the Handle method to
// process any errors from Handlers.
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

	attrs := map[string]any{}
	record.Attrs(func(attr slog.Attr) bool {
		handler.appendAttr(attr, attrs, entry)
		return true
	})
	entry.attributes = attrs

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
func (handler Handler) appendAttr(attr slog.Attr, attrs map[string]any, entry *entry) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}

	kind := attr.Value.Kind()
	if kind == slog.KindGroup {
		groupAttrs := attr.Value.Group()
		if len(groupAttrs) == 0 {
			return
		}
		// If the key is non-empty, write it out and indent the rest of the attrs.
		// Otherwise, inline the attrs.
		if attr.Key == "" { // Inline
			for _, attr := range groupAttrs {
				handler.appendAttr(attr, attrs, entry)
			}
		} else {
			groupAttr := map[string]any{}
			for _, attr := range groupAttrs {
				handler.appendAttr(attr, groupAttr, nil) // Ignore the special keys if they're in a group
			}
			attrs[attr.Key] = groupAttr
		}
	}

	if entry != nil && attr.Key == PublicMessageKey {
		entry.publicMessage = attr.Value.String()
	} else if entry != nil && attr.Key == UserIDKey {
		intValue, ok := attr.Value.Any().(int)
		if ok {
			entry.userID = intValue
		} else {
			typeOf := reflect.TypeOf(attr.Value.Any())
			fmt.Printf("warning: userID property in log statement was not an int so has been ignored. type: %v", typeOf) // TODO: can the logger call itself?
		}
		attrs[attr.Key] = attr.Value.Any() // Also store the value in the attributes so it's preserved if the user is deleted
	} else {
		attrs[attr.Key] = attr.Value.Any()
	}
}

// func (handler *Handler) appendAttr(buffer []byte, attr slog.Attr) []byte {
// 	attr.Value = attr.Value.Resolve()
// 	if attr.Equal(slog.Attr{}) {
// 		return buffer
// 	}

// 	switch attr.Value.Kind() {
// 	case slog.KindGroup:
// 		// TODO: implement
// 	default:
// 		if len(buffer) != 0 {
// 			buffer = append(buffer, " "...)
// 		}
// 		buffer = append(buffer, attr.Key...)
// 		buffer = append(buffer, "="...)
// 		buffer = append(buffer, attr.Value.String()...)
// 	}

// 	return buffer
// }

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (handler Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handler.textHandler = handler.textHandler.WithAttrs(attrs)
	panic("not implemented")
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
	handler.textHandler = handler.textHandler.WithGroup(name)
	panic("not implemented")
	return handler
}
