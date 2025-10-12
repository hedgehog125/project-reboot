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

var ErrWrapperStartTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeStartTx).
	SetChild(ErrWrapperDatabase)
var ErrWrapperCommitTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCommitTx).
	SetChild(ErrWrapperDatabase)
var ErrWrapperCallback = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeCallback)
var ErrWrapperWithTx = common.NewErrorWrapper(common.ErrTypeDbCommon, ErrTypeWithTx)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeDbCommon).
	SetChild(common.ErrWrapperDatabase)

var ErrUnexpectedTransaction = ErrWrapperStartTx.Wrap(
	ErrWrapperDatabase.Wrap(
		errors.New("found transaction in context. nested transactions are not supported"),
	),
)
