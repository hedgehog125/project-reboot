package core

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeSendActiveSessionReminders = "send active session reminders"
	ErrTypeDeleteExpiredSessions      = "delete expired sessions"
	ErrTypeEncrypt                    = "encrypt"
	ErrTypeDecrypt                    = "decrypt"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrWrapperInvalidData = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeInvalidData)
var ErrWrapperCreateCipher = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeInvalidData)

var ErrIncorrectPassword = common.NewErrorWithCategories("incorrect password", common.ErrTypeCore)

var ErrWrapperSendActiveSessionReminders = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeSendActiveSessionReminders)
var ErrWrapperDeleteExpiredSessions = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeDeleteExpiredSessions)

// These functions don't categorize their errors
var ErrWrapperEncrypt = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeEncrypt)
var ErrWrapperDecrypt = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeDecrypt)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeCore).
	SetChild(common.ErrWrapperDatabase)
