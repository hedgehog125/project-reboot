package messengers

import (
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/util"
)

type Messenger interface {
	Send(message Message) error
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

type MessengerGroup interface {
	SendUsingAll(message Message) []error
}

type messengerGroup struct {
	messengers []Messenger
}

func (group *messengerGroup) SendUsingAll(message Message) []error {
	errChan := make(chan util.ErrWithIndex, 3)
	for i, messenger := range group.messengers {
		go func() {
			err := messenger.Send(message)
			errChan <- util.ErrWithIndex{
				Err:   err,
				Index: i,
			}
		}()
	}
	errs := make([]error, len(group.messengers))
	for range len(group.messengers) {
		errInfo := <-errChan
		errs[errInfo.Index] = errInfo.Err
	}

	return errs
}

func NewMessengerGroup(messengers []Messenger) MessengerGroup {
	return &messengerGroup{
		messengers: messengers,
	}
}
