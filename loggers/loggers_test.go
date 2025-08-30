package loggers_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/logentry"
	"github.com/hedgehog125/project-reboot/loggers"
	"github.com/jonboulle/clockwork"
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
	UserID        int
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
		require.Equal(t, expected.UserID, entry.UserID)
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

func TestLogger_SpecialAttributes(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	userIDs := []int{}
	for i := range 2 {
		userOb := db.Client().User.Create().SetUsername(fmt.Sprintf("user%v", i+1)).
			SetContent([]byte{1}).SetFileName("file.zip").SetMime("application/zip").
			SetNonce([]byte{1}).SetKeySalt([]byte{1}).
			SetHashTime(0).SetHashMemory(0).SetHashThreads(0).
			SaveX(t.Context())
		userIDs = append(userIDs, userOb.ID)
	}

	defer db.Shutdown()
	app := &common.App{
		Database: db,
		Env:      testcommon.DefaultEnv(),
		Clock:    clockwork.NewRealClock(),
	}
	logger := NewLogger(app)
	app.Logger = logger
	logger.Start()

	logger.Info("deleted expired sessions", loggers.UserIDKey, userIDs[0])
	logger.Info(
		"public message nobody will be sent",
		loggers.PublicMessageKey, "public version of \"public message nobody will be sent\"",
	)
	logger.Info(
		"updated password",
		loggers.UserIDKey, userIDs[1],
		loggers.PublicMessageKey,
		"your password was updated",
		"hiddenData",
		"shh",
	)

	time.Sleep(5 * time.Millisecond)
	logger.Shutdown()
	logger.AssertWritten(t, []ExpectedEntry{
		{
			Message: "deleted expired sessions",
			Level:   int(slog.LevelInfo),
			UserID:  1,
			Attributes: map[string]any{
				"userID": userIDs[0],
			},
		},
		{
			Message:       "public message nobody will be sent",
			Level:         int(slog.LevelInfo),
			PublicMessage: "public version of \"public message nobody will be sent\"",
		},
		{
			Message:       "updated password",
			Level:         int(slog.LevelInfo),
			PublicMessage: "your password was updated",
			UserID:        userIDs[1],
			Attributes: map[string]any{
				"userID":     userIDs[1],
				"hiddenData": "shh",
			},
		},
	})
}

func TestLogger_RetriesBulkCreateIndividually(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	defer db.Shutdown()

	var successfulCreateCounter atomic.Int64
	var createAttemptCounter atomic.Int64
	var pendingMutations = common.MutexValue[map[*ent.Tx]int64]{
		Value: map[*ent.Tx]int64{},
	}
	db.AddStartTxHook(func(tx *ent.Tx) error {
		tx.LogEntry.Use(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
				if mutation.Op().Is(ent.OpCreate) {
					pendingMutations.Mutex.Lock()
					pendingMutations.Value[tx]++
					pendingMutations.Mutex.Unlock()
					if createAttemptCounter.Add(1) <= 1 {
						return nil, errors.New("temporary but unretryable error")
					}
				}
				return next.Mutate(ctx, mutation)
			})
		})
		tx.OnCommit(func(committer ent.Committer) ent.Committer {
			return ent.CommitFunc(
				func(ctx context.Context, tx *ent.Tx) error {
					stdErr := committer.Commit(ctx, tx)
					if stdErr != nil {
						return stdErr
					}

					pendingMutations.Mutex.Lock()
					successfulCreateCounter.Add(pendingMutations.Value[tx])
					delete(pendingMutations.Value, tx)
					pendingMutations.Mutex.Unlock()
					return nil
				},
			)
		})
		return nil
	})

	app := &common.App{
		Database: db,
		Env:      testcommon.DefaultEnv(),
		Clock:    clockwork.NewRealClock(),
	}
	logger := NewLogger(app)
	app.Logger = logger
	logger.Start()

	logger.Info("doing something")
	logger.Info("doing something else")

	time.Sleep(5 * time.Millisecond)
	logger.Shutdown()
	logger.AssertWritten(t, []ExpectedEntry{
		{
			Message: "doing something",
			Level:   int(slog.LevelInfo),
		},
		{
			Message: "doing something else",
			Level:   int(slog.LevelInfo),
		},
		{
			Message: "bulk log write failed but the individual fallback writes all succeeded, so the writes took longer than they should have",
			Level:   int(slog.LevelWarn),
			Attributes: map[string]any{
				"error": map[string]any{
					"categories": []any{
						"callback",
						"WithTx",
						"db common [package]",
					},
					"debugValues": []any{
						map[string]any{
							"Value":   []any{},
							"message": "no previous errors",
							"name":    "previous retry errors (WithRetries)",
						},
						map[string]any{
							"Value":   nil,
							"message": "max retries: 0, base backoff: 0s, backoff multiplier: 0",
							"name":    "retries reset by WithRetries from...",
						},
					},
					"errDuplicatesCategory":  false,
					"error":                  "db common [package] error: WithTx error: callback error: temporary but unretryable error",
					"innerError":             "temporary but unretryable error",
					"maxRetries":             0,
					"retryBackoffBase":       0,
					"retryBackoffMultiplier": 0,
				},
			},
		},
	})
	// Should do a bulk create that fails, retry that with 2 individual creates and then do another bulk create to store the warning
	require.Equal(t, int64(3), successfulCreateCounter.Load())
}
