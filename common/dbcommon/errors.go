package dbcommon

import (
	"errors"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeWithTx = "WithTx"
	// Lower level
	ErrTypeStartTx  = "start transaction"
	ErrTypeCommitTx = "commit transaction"
	ErrTypeCallback = "callback"
)

// Not used by this package. Return this error when you need to cancel the transaction and don't have an error
var ErrCancelTransaction = common.NewErrorWithCategories("cancel transaction")

var ErrWrapperStartTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeStartTx)
var ErrWrapperCommitTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCommitTx)
var ErrWrapperCallback = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCallback)
var ErrWrapperWithTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeWithTx)

var ErrUnexpectedTransaction = ErrWrapperStartTx.Wrap(
	errors.New("found transaction in context. nested transactions are not supported"),
)
