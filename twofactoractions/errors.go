package twofactoractions

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeCreate       = "create"
	ErrTypeConfirm      = "confirm"
	ErrTypeEncode       = "encode"
	ErrTypeDecodeAction = "decode action" // From Action.Decode() method
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

// TODO: test this
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeTwoFactorAction).SetChild(common.ErrWrapperDatabase)
var ErrWrapperInvalidData = common.NewErrorWrapper(ErrTypeInvalidData, common.ErrTypeTwoFactorAction)

var ErrWrapperDecodeAction = common.NewErrorWrapper(
	ErrTypeDecodeAction, common.ErrTypeTwoFactorAction,
)
