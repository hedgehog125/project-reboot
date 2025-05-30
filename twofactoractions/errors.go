package twofactoractions

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeCreate  = "create"
	ErrTypeConfirm = "confirm"
	ErrTypeEncode  = "encode"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrNotFound = common.NewErrorWithCategories(
	"no action with given ID", common.ErrTypeTwoFactorAction,
)
var ErrWrongCode = common.NewErrorWithCategories(
	"wrong 2FA code", common.ErrTypeTwoFactorAction,
)
var ErrExpired = common.NewErrorWithCategories(
	"action has expired", common.ErrTypeTwoFactorAction,
)
var ErrUnknownActionType = common.NewErrorWithCategories(
	"unknown action type", common.ErrTypeTwoFactorAction,
)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeDatabase, common.ErrTypeTwoFactorAction)
var ErrWrapperInvalidData = common.NewErrorWrapper(ErrTypeInvalidData, common.ErrTypeTwoFactorAction)
