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

		if expected.Attributes == nil {
			expected.Attributes = map[string]any{}
		}
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

func TestLogger_WithAttrs_and_WithGroup(t *testing.T) {
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

	logger.With("requestID", "request-1").Info("created user", "userID", 1)
	logger.WithGroup("group").Info("created user", "userID", 2)
	logger.With(
		"requestID", "request-2",
	).WithGroup("user_data").Info("user session started", "username", "alice")
	logger.WithGroup("http").WithGroup("request").Info("incoming request 1", "method", "GET", "path", "/api/v1/users")
	logger.WithGroup("http").Info("incoming request 2", slog.Group("request", "method", "GET", "path", "/api/v1/users"))

	jobLogger := logger.With(slog.Int("jobID", 3))
	jobCallbackLogger := jobLogger.WithGroup("callback")
	jobCallbackLogger.With("callbackValue1", 4).Debug("doing something")
	jobLogger.With("timeTaken", time.Second).Debug("job did something")
	jobCallbackLogger.Debug("doing something else", "callbackValue2", 5)

	logger.Warn("simple warning")

	time.Sleep(5 * time.Millisecond)
	logger.Shutdown()
	logger.AssertWritten(t, []ExpectedEntry{
		{
			Message: "created user",
			Level:   int(slog.LevelInfo),
			Attributes: map[string]any{
				"requestID": "request-1",
				"userID":    1,
			},
		},
		{
			Message: "created user",
			Level:   int(slog.LevelInfo),
			Attributes: map[string]any{
				"group": map[string]any{
					"userID": 2,
				},
			},
		},
		{
			Message: "user session started",
			Level:   int(slog.LevelInfo),
			Attributes: map[string]any{
				"requestID": "request-2",
				"user_data": map[string]any{
					"username": "alice",
				},
			},
		},
		{
			Message: "incoming request 1",
			Level:   int(slog.LevelInfo),
			Attributes: map[string]any{
				"http": map[string]any{
					"request": map[string]any{
						"method": "GET",
						"path":   "/api/v1/users",
					},
				},
			},
		},
		{
			Message: "incoming request 2",
			Level:   int(slog.LevelInfo),
			Attributes: map[string]any{
				"http": map[string]any{
					"request": map[string]any{
						"method": "GET",
						"path":   "/api/v1/users",
					},
				},
			},
		},
		{
			Message: "doing something",
			Level:   int(slog.LevelDebug),
			Attributes: map[string]any{
				"jobID": 3,
				"callback": map[string]any{
					"callbackValue1": 4,
				},
			},
		},
		{
			Message: "job did something",
			Level:   int(slog.LevelDebug),
			Attributes: map[string]any{
				"jobID":     3,
				"timeTaken": time.Second,
			},
		},
		{
			Message: "doing something else",
			Level:   int(slog.LevelDebug),
			Attributes: map[string]any{
				"jobID": 3,
				"callback": map[string]any{
					"callbackValue2": 5,
				},
			},
		},
		{
			Message: "simple warning",
			Level:   int(slog.LevelWarn),
		},
	})
}
