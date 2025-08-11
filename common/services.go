package common

/*
The core principal is to abstract just enough that:
* The service can be mocked to some extent (although I don't think this is really necessary for the database)
* The service can be used in simplified ways for testing. e.g a test can use a different job registry with a real implementation
*/

import (
	"context"
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

	UNLOCK_TIME time.Duration
	// TODO: implement
	AUTH_CODE_VALID_FOR time.Duration

	PASSWORD_HASH_SETTINGS *PasswordHashSettings

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
	State            *State
	TwoFactorActions TwoFactorActionService
	Messengers       MessengerService // TODO: does this still need to be a service?
	Database         DatabaseService
	Server           ServerService
	Jobs             JobService
	Scheduler        SchedulerService
}

type MessengerService interface {
	// Note: this atomically queues jobs to send the messages
	SendUsingAll(message *Message, ctx context.Context) *Error
}
type MessageType string

const (
	MessageTest     = "test"
	MessageRegular  = "regular"
	MessageLogin    = "login"
	MessageReset    = "reset"
	MessageLock     = "lock"
	MessageSelfLock = "selfLock"
	Message2FA      = "2FA"
)

type Message struct {
	Type  MessageType
	User  *ent.User
	Code  string
	Until time.Time
}

type DatabaseService interface {
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
		data any,
		ctx context.Context,
	) (uuid.UUID, *Error)
	EnqueueEncoded(
		versionedType string,
		encodedData string,
		ctx context.Context,
	) (uuid.UUID, *Error)
	WaitForJobs()
	Encode(versionedType string, data any) (string, *Error)
}
type TwoFactorActionService interface {
	Create(
		versionedType string,
		expiresAt time.Time,
		data any,
		ctx context.Context,
	) (uuid.UUID, string, *Error)
	Confirm(actionID uuid.UUID, code string, ctx context.Context) (uuid.UUID, *Error)
}

type SchedulerService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}
