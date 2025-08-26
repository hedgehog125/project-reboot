package loggers_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/logentry"
	"github.com/hedgehog125/project-reboot/loggers"
	"github.com/stretchr/testify/require"
)

type Logger struct {
	*slog.Logger
	Handler loggers.Handler
}

func NewLogger(app *common.App) *Logger {
	handler := loggers.NewHandler(slog.LevelDebug, true, true, app)
	return &Logger{
		Logger:  slog.New(handler),
		Handler: handler,
	}
}

func (service *Logger) Start() {
	go service.Handler.Listen()
}
func (service *Logger) Shutdown() {
	service.Handler.Shutdown()
}

type ExpectedEntry struct {
	Message       string
	PublicMessage string
	Level         int
	Attributes    map[string]any
	// UserID should be asserted by its attribute
}

func (service *Logger) AssertWritten(t *testing.T, expectedEntries []ExpectedEntry) {
	client := service.Handler.App.Database.Client()
	entries := client.LogEntry.Query().Order(ent.Asc(logentry.FieldTime)).AllX(t.Context())
	require.Len(t, entries, len(expectedEntries))
	for i, entry := range entries {
		expected := expectedEntries[i]
		require.Equal(t, expected.Message, entry.Message)
		require.Equal(t, expected.PublicMessage, entry.PublicMessage)
		require.Equal(t, expected.Level, entry.Level)
		testcommon.AssertJSONEqual(t, expected.Attributes, entry.Attributes)
	}
}

func TestLogger_SavesToDatabase(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	defer db.Shutdown()
	app := &common.App{
		Database: db,
		Env:      testcommon.DefaultEnv(),
	}
	logger := NewLogger(app)
	app.Logger = logger
	logger.Start()

	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warning")
	logger.Error("error")

	time.Sleep(5 * time.Millisecond)
	logger.Shutdown()
	logger.AssertWritten(t, []ExpectedEntry{
		{
			Message: "debug",
			Level:   int(slog.LevelDebug),
		},
		{
			Message: "info",
			Level:   int(slog.LevelInfo),
		},
		{
			Message: "warning",
			Level:   int(slog.LevelWarn),
		},
		{
			Message: "error",
			Level:   int(slog.LevelError),
		},
	})
}

// func TestLogger_WithAttrs_and_WithGroup(t *testing.T) {
// 	t.Parallel()
// 	db := testcommon.CreateDB()
// 	defer db.Shutdown()
// 	app := &common.App{
// 		Database: db,
// 		Env:      testcommon.DefaultEnv(),
// 	}
// 	logger := NewLogger(app)
// }
