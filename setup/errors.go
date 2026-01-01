package setup

import "github.com/NicoClack/cryptic-stash/common"

const (
	ErrTypeGenerateAdminSetupConstants = "generate admin setup constants"
)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeSetup).
	SetChild(common.ErrWrapperDatabase)
var ErrWrapperTotp = common.NewErrorWrapper(common.ErrTypeSetup, common.ErrTypeTotp)

var ErrWrapperGenerateAdminSetupConstants = common.NewErrorWrapper(
	common.ErrTypeSetup, ErrTypeGenerateAdminSetupConstants,
)
