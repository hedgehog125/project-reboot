package messagers

import "github.com/hedgehog125/project-reboot/ent"

type Messager interface {
	SendBatch(messages []Message) error
}

type MessageType string

const (
	MessageRegular = "regular"
	MessageLogin   = "login"
)

type Message struct {
	Type MessageType
	// Won't include sensitive properties like Content
	User *ent.User
}
