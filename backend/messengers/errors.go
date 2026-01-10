package messengers

import (
	"github.com/NicoClack/cryptic-stash/backend/common"
)

const (
	ErrTypeSend         = "send"
	ErrTypeSendUsingAll = "send using all"
	ErrTypeSendBulk     = "send bulk"
	// Lower level
	ErrTypeFormatMessage = "format message"
	ErrTypePrepare       = "prepare"
	ErrTypeEnqueueJob    = "enqueue job"
)

var ErrWrapperPrepare = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypePrepare,
)
var ErrWrapperEnqueueJob = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypeEnqueueJob,
)
var ErrWrapperSend = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeSend)
var ErrWrapperSendUsingAll = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeSendUsingAll)
var ErrWrapperSendBulk = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeSendBulk)

var ErrWrapperDecodeOptions = common.NewErrorWrapper(common.ErrTypeMessengers, "decode options")
var ErrMessengerDisabledForUser = common.NewErrorWithCategories(
	"messenger type disabled for user", common.ErrTypeMessengers,
)
var ErrUnknownMessengerType = common.NewErrorWithCategories(
	"unknown messenger type", common.ErrTypeMessengers,
)

var ErrWrapperFormat = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeFormatMessage)
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).
	SetChild(common.ErrWrapperDatabase)

// Job errors
// Messengers wrap the Handle method of the job. If there's an error from the outer handler, we use this wrapper
var ErrWrapperHandlerWrapper = common.NewErrorWrapper("handler wrapper")
