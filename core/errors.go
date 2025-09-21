package core

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeEncrypt                    = "encrypt"
	ErrTypeDecrypt                    = "decrypt"
	ErrTypeSendActiveSessionReminders = "send active session reminders"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrWrapperInvalidData = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeInvalidData)
var ErrWrapperCreateCipher = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeInvalidData)

var ErrIncorrectPassword = common.NewErrorWithCategories("incorrect password", common.ErrTypeCore)

// These functions don't categorize their errors
var ErrWrapperEncrypt = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeEncrypt)
var ErrWrapperDecrypt = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeDecrypt)

var ErrWrapperSendActiveSessionReminders = common.NewErrorWrapper(common.ErrTypeCore, ErrTypeSendActiveSessionReminders)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeCore).
	SetChild(common.ErrWrapperDatabase)
