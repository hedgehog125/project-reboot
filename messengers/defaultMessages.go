package messengers

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
)

var defaultMessageMap = map[common.MessageType]func(message common.Message) string{
	common.MessageLogin: func(message common.Message) string {
		return fmt.Sprintf("Login attempt")
	},
	common.MessageTest: func(message common.Message) string {
		return "Test message"
	},
}

// For messengers like SMS where the messages should be as short as possible with no formatting
func formatDefaultMessage(message common.Message) (string, error) {
	formatter, ok := defaultMessageMap[message.Type]
	if !ok {
		return "", fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type)
	}

	return formatter(message), nil
}
