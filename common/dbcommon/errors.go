package dbcommon

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeWithTx = "WithTx"
	// Lower level
	ErrTypeStartTx  = "start transaction"
	ErrTypeCommitTx = "commit transaction"
	ErrTypeCallback = "callback"
)

var ErrWrapperStartTx = common.NewErrorWrapper(ErrTypeStartTx, common.ErrTypeDbCommon)
var ErrWrapperCommitTx = common.NewErrorWrapper(ErrTypeCommitTx, common.ErrTypeDbCommon)
var ErrWrapperCallback = common.NewErrorWrapper(ErrTypeCallback, common.ErrTypeDbCommon)
