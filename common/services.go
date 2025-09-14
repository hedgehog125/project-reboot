package common

/*
The core principal is to abstract just enough that:
* The service can be mocked to some extent (although I don't think this is really necessary for the database)
* The service can be used in simplified ways for testing. e.g a test can use a different job registry with a real implementation
*/

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/jonboulle/clockwork"
)

type Env struct {
	IS_DEV                        bool
	PORT                          int
	MOUNT_PATH                    string
	PROXY_ORIGINAL_IP_HEADER_NAME string

	JOB_POLL_INTERVAL    time.Duration
	MAX_TOTAL_JOB_WEIGHT int

	UNLOCK_TIME         time.Duration
	AUTH_CODE_VALID_FOR time.Duration
	// Once used, how much longer the auth code remains valid for
	USED_AUTH_CODE_VALID_FOR time.Duration
	PASSWORD_HASH_SETTINGS   *PasswordHashSettings

	LOG_STORE_INTERVAL time.Duration
	ADMIN_USERNAME     string
	// How long the server should wait for messengers to succeed before crashing the server to send the message
	// Note: this time will be exceeded as it's a simple check when the job succeeds and doesn't take into account when the next retry is
	// Note: currently all of the successfully prepared messages must succeed for a crash to be avoided
	ADMIN_MESSAGE_TIMEOUT time.Duration

	DISCORD_TOKEN  string
	SENDGRID_TOKEN string // TODO: implement
}
type PasswordHashSettings struct {
	Time   uint32
	Memory uint32
	// Note: this affects the hash produced
	Threads uint8
}

type State struct {
	AdminCode chan []byte
}
type App struct {
	Env              *Env
	Clock            clockwork.Clock
	Logger           LoggerService
	ShutdownService  ShutdownService
	Database         DatabaseService
	State            *State
	TwoFactorActions TwoFactorActionService
	Messengers       MessengerService
	Server           ServerService
	Jobs             JobService
	Scheduler        SchedulerService
}

// If reason is "", the server will exit with a 0 exit code
func (app *App) Shutdown(reason string) {
	go app.ShutdownService.Shutdown(reason)
}

type HasDefaultLogger interface {
	DefaultLogger() Logger
}

type MessengerService interface {
	// The error map is more like warnings about why specific messengers failed to prepare, they are logged already so you might just want to ignore them
	//
	// But check the second *Error value first because you should fail the transaction if it's not nil
	//
	// Note: the number of successfully queued messages (the int return value) might not be 0 if some messages were queued before a non-messenger specific error occurred
	SendUsingAll(message *Message, ctx context.Context) (int, map[string]*Error, *Error)
}
type MessageType string

const (
	MessageTest       = "test"
	MessageAdminError = "adminError"
	MessageRegular    = "regular"
	MessageLogin      = "login"
	MessageDownload   = "download"
	MessageUserUpdate = "userUpdate"
	MessageLock       = "lock"
	MessageSelfLock   = "selfLock"
	Message2FA        = "2FA"
)

type Message struct {
	Type MessageType
	User *ent.User
	Code string
	Time time.Time
}

type Logger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Enabled(ctx context.Context, level slog.Level) bool
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Log(ctx context.Context, level slog.Level, msg string, args ...any)
	LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
}
type LoggerService interface {
	Logger
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}

// When in a context passed to a logger.Error call, the server will deliberately crash to notify the admin as opposed to sending a message
type AdminNotificationFallbackKey struct{}

// Used to store a logger override in a context
type LoggerKey struct{}

type ShutdownService interface {
	// Note: this blocks until shutdown is complete, crashes should usually call this in a separate Goroutine
	//
	// If reason is "", the server will exit with a 0 exit code
	Shutdown(reason string)
}

type DatabaseService interface {
	HasDefaultLogger
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
	Client() *ent.Client
	ReadTx(ctx context.Context) (*ent.Tx, error)
	WriteTx(ctx context.Context) (*ent.Tx, error)
}

type ServerService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}

type JobService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
	Enqueue(
		versionedType string,
		body any,
		ctx context.Context,
	) (*ent.Job, *Error)
	EnqueueEncoded(
		versionedType string,
		encodedBody json.RawMessage,
		ctx context.Context,
	) (*ent.Job, *Error)
	EnqueueWithModifier(
		versionedType string,
		body any,
		modifications func(*ent.JobCreate) *ent.JobCreate,
		ctx context.Context,
	) (*ent.Job, *Error)
	EnqueueEncodedWithModifier(
		versionedType string,
		encodedBody json.RawMessage,
		modifications func(*ent.JobCreate) *ent.JobCreate,
		ctx context.Context,
	) (*ent.Job, *Error)
	WaitForJobs()
	Encode(versionedType string, body any) (json.RawMessage, *Error)
}
type TwoFactorActionService interface {
	Create(
		versionedType string,
		expiresAt time.Time,
		body any,
		ctx context.Context,
	) (*ent.TwoFactorAction, string, *Error)
	Confirm(actionID uuid.UUID, code string, ctx context.Context) (*ent.Job, *Error)
}

type SchedulerService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}
