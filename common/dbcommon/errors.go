package dbcommon

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeWithTx = "WithTx"
	// Lower level
	ErrTypeStartTx  = "start transaction"
	ErrTypeCommitTx = "commit transaction"
	ErrTypeCallback = "callback"
)

var ErrWrapperStartTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeStartTx)
var ErrWrapperCommitTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCommitTx)
var ErrWrapperCallback = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCallback)
var ErrWrapperWithTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeWithTx)
