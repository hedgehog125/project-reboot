package messengers

import "github.com/hedgehog125/project-reboot/ent"

type Messenger interface {
	SendBatch(messages []Message) error
}

type MessageType string

const (
	MessageTest    = "test"
	MessageRegular = "regular"
	MessageLogin   = "login"
)

type Message struct {
	Type MessageType
	// Won't include sensitive properties like Content
	User *ent.User
}
