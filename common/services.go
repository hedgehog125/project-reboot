package common

import (
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/jonboulle/clockwork"
)

type Env struct {
	PORT                          int
	MOUNT_PATH                    string
	PROXY_ORIGINAL_IP_HEADER_NAME string
	UNLOCK_TIME                   int64 // In seconds
	// TODO: implement
	AUTH_CODE_VALID_FOR int64 // In seconds

	DISCORD_TOKEN  string
	SENDGRID_TOKEN string // TODO: implement
}
type State struct {
	AdminCode chan []byte
}
type App struct {
	Env       *Env
	Clock     clockwork.Clock
	State     *State
	Messenger MessengerService
	Database  DatabaseService
	Server    ServerService
	Scheduler SchedulerService
}

type MessengerService interface {
	SendUsingAll(message Message) []error
}
type MessageType string

const (
	MessageTest    = "test"
	MessageRegular = "regular"
	MessageLogin   = "login"
)

type Message struct {
	Type MessageType
	User *MessageUserInfo
}

// The info about the user provided to a Messenger
type MessageUserInfo struct {
	Username       string
	AlertDiscordId string
	AlertEmail     string
}

type DatabaseService interface {
	Client() *ent.Client
	Shutdown() // Should log warning rather than return an error
}

type ServerService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}

type SchedulerService interface {
	Start()    // Should fatalf rather than returning an error
	Shutdown() // Should log warning rather than return an error
}
