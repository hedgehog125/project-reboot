package twofactoractions

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeCreate   = "create"
	ErrTypeConfirm  = "confirm"
	ErrTypeNotFound = "not found"
)

var ErrNoTxInContext = common.ErrNoTxInContext.AddCategory(common.ErrTypeTwoFactorAction)
var ErrNotFound = common.NewErrorWithCategories(
	common.ErrTypeTwoFactorAction, "no action with given ID",
)
var ErrWrongCode = common.NewErrorWithCategories(
	common.ErrTypeTwoFactorAction, "wrong 2FA code",
)
var ErrExpired = common.NewErrorWithCategories(
	common.ErrTypeTwoFactorAction, "action has expired",
)

var ErrWrapperCreate = common.NewErrorWrapper(common.ErrTypeTwoFactorAction, ErrTypeCreate)
var ErrWrapperConfirm = common.NewErrorWrapper(common.ErrTypeTwoFactorAction, ErrTypeConfirm)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeTwoFactorAction).SetChild(common.ErrWrapperDatabase)
