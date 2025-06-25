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
	common.Message2FA: func(message common.Message) string {
		return fmt.Sprintf("2FA code: %s", message.Code)
	},
	common.MessageSelfLock: func(message common.Message) string {
		return fmt.Sprintf("You have locked your account until %s", message.Until.Format("2006-01-02 15:04:05"))
	},
}

// For messengers like SMS where the messages should be as short as possible with no formatting
func formatDefaultMessage(message common.Message) (string, *common.Error) {
	formatter, ok := defaultMessageMap[message.Type]
	if !ok {
		return "", ErrWrapperFormat.Wrap(
			fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type),
		)
	}

	return formatter(message), nil
}
