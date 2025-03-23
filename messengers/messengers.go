package messengers

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
)

type Messenger interface {
	Id() string
	Send(message common.Message) error
}

// For messengers like SMS where the messages should be as short as possible with no formatting
func formatDefaultMessage(message common.Message) (string, error) {
	switch message.Type {
	case common.MessageLogin:
		return fmt.Sprintf("Login attempt"), nil
	case common.MessageTest:
		return "Test message", nil
	}
	return "", fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type)
}
