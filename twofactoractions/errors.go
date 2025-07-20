package twofactoractions

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeCreate   = "create"
	ErrTypeConfirm  = "confirm"
	ErrTypeNotFound = "not found"
)

var ErrNoTxInContext = common.ErrNoTxInContext.AddCategory(common.ErrTypeTwoFactorAction)
var ErrNotFound = common.NewErrorWithCategories(
	"no action with given ID", common.ErrTypeTwoFactorAction,
)
var ErrWrongCode = common.NewErrorWithCategories(
	"wrong 2FA code", common.ErrTypeTwoFactorAction,
)
var ErrExpired = common.NewErrorWithCategories(
	"action has expired", common.ErrTypeTwoFactorAction,
)

var ErrWrapperCreate = common.NewErrorWrapper(ErrTypeCreate, common.ErrTypeTwoFactorAction)
var ErrWrapperConfirm = common.NewErrorWrapper(ErrTypeConfirm, common.ErrTypeTwoFactorAction)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeTwoFactorAction).SetChild(common.ErrWrapperDatabase)
