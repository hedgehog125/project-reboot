package twofactoractions

import "github.com/NicoClack/cryptic-stash/backend/common"

const (
	ErrTypeCreate               = "create"
	ErrTypeConfirm              = "confirm"
	ErrTypeDeleteExpiredActions = "delete expired actions"
	ErrTypeNotFound             = "not found"
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

var ErrWrapperCreate = common.NewErrorWrapper(common.ErrTypeTwoFactorAction, ErrTypeCreate)
var ErrWrapperConfirm = common.NewErrorWrapper(common.ErrTypeTwoFactorAction, ErrTypeConfirm)
var ErrWrapperDeleteExpiredActions = common.NewErrorWrapper(common.ErrTypeTwoFactorAction, ErrTypeDeleteExpiredActions)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeTwoFactorAction).SetChild(common.ErrWrapperDatabase)
