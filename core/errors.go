package core

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeEncrypt = "encrypt"
	ErrTypeDecrypt = "decrypt"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrWrapperInvalidData = common.NewErrorWrapper(ErrTypeInvalidData, common.ErrTypeCore)
var ErrWrapperCreateCipher = common.NewErrorWrapper(ErrTypeInvalidData, common.ErrTypeCore)

var ErrIncorrectPassword = common.NewErrorWithCategories("incorrect password")

// These functions don't categorize their errors
var ErrWrapperEncrypt = common.NewErrorWrapper(ErrTypeEncrypt, common.ErrTypeCore)
var ErrWrapperDecrypt = common.NewErrorWrapper(ErrTypeDecrypt, common.ErrTypeCore)
