package messengers

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
)

func getLoginAttemptMessageBody(message *common.Message) string {
	explanation := "If you're reading this, it most likely wasn't you that's logging in! Please self-lock your user ASAP and contact your admin as they can make the lock permanent if you want. If you want to be able to safely unlock your user, you should update your password with help from your admin."
	return fmt.Sprintf(
		"%v\n\nIF YOU DO NOTHING, we'll assume you're locked out and will ALLOW THE USER TO LOG IN after %v UTC.",
		explanation,
		message.Time.Format("2006-01-02 15:04:05"),
	)
}

var defaultMessageMap = map[common.MessageType]func(message *common.Message) string{
	common.MessageUserUpdate: func(message *common.Message) string {
		return "Your account password and/or file have been updated by your admin."
	},
	common.MessageLogin: func(message *common.Message) string {
		return "LOGIN ATTEMPT! " + getLoginAttemptMessageBody(message)
	},
	common.MessageActiveSessionReminder: func(message *common.Message) string {
		return "REMINDER: YOU HAVE A PENDING LOGIN ATTEMPT! " + getLoginAttemptMessageBody(message)
	},
	common.MessageDownload: func(message *common.Message) string {
		return "Your data has been downloaded. If this wasn't you, please rotate your 2FA backup codes immediately and contact your admin!"
	},
	common.MessageTest: func(message *common.Message) string {
		return "If you're reading this message, it means your updated contacts are working."
	},
	common.Message2FA: func(message *common.Message) string {
		return fmt.Sprintf("2FA code: %s", message.Code)
	},
	common.MessageSelfLock: func(message *common.Message) string {
		return fmt.Sprintf("You have locked your account until %s", message.Time.Format("2006-01-02 15:04:05"))
	},
	common.MessageAdminError: func(message *common.Message) string {
		return "[Admin] An error has occurred! Please investigate the logs and possibly create an issue at https://github.com/hedgehog125/project-reboot/issues as this might be reducing security"
	},
}

// For messengers like SMS where the messages should be as short as possible with no formatting
func FormatDefaultMessage(message *common.Message) (string, *common.Error) {
	formatter, ok := defaultMessageMap[message.Type]
	if !ok {
		return "", ErrWrapperFormat.Wrap(
			fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type),
		)
	}

	return formatter(message), nil
}
