package setup

import "github.com/NicoClack/cryptic-stash/backend/common"

const (
	ErrTypeGenerateAdminSetupConstants = "generate admin setup constants"
	ErrTypeGetStatus                   = "get status"
	// Lower level
	ErrTypeCheckAdminHasMessengers = "check admin has messengers"
)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeSetup).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperTotp = common.NewErrorWrapper(common.ErrTypeSetup, common.ErrTypeTotp)

var ErrWrapperGenerateAdminSetupConstants = common.NewErrorWrapper(
	common.ErrTypeSetup, ErrTypeGenerateAdminSetupConstants,
)
var ErrWrapperGetStatus = common.NewErrorWrapper(
	common.ErrTypeSetup, ErrTypeGetStatus,
)
var ErrWrapperCheckAdminHasMessengers = common.NewErrorWrapper(
	common.ErrTypeSetup, ErrTypeCheckAdminHasMessengers,
)
