package common

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
	Messenger        MessengerService // TODO: does this still need to be a service?
	Database         DatabaseService
	Server           ServerService
	Jobs             JobService
	Scheduler        SchedulerService
}

type MessengerService interface {
	IDs() []string
	// Note: this should be atomic and call the messengers in the background (usually via jobs)
	SendUsingAll(message Message) *Error
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
	User  *UserContacts
	Code  string
	Until time.Time
}

// The info about the user provided to a Messenger
type UserContacts struct {
	Username       string
	AlertDiscordId string
	AlertEmail     string
}

type DatabaseService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
	Client() *ent.Client
	Tx(ctx context.Context) (*ent.Tx, error)
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
