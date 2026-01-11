package messengers

import (
	"github.com/NicoClack/cryptic-stash/backend/common"
)

const (
	ErrTypeSend            = "send"
	ErrTypeSendUsingAll    = "send using all"
	ErrTypeSendBulk        = "send bulk"
	ErrTypeEnableMessenger = "enable messenger"
	// Lower level
	ErrTypeDecodeOptions = "decode options"
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
var ErrWrapperEnableMessenger = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeEnableMessenger)

var ErrWrapperDecodeOptions = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeDecodeOptions)
var ErrMessengerDisabledForUser = common.NewErrorWithCategories(
	"messenger type disabled for user", common.ErrTypeMessengers,
)
var ErrUnknownMessengerType = common.NewErrorWithCategories(
	"unknown messenger type", common.ErrTypeMessengers,
)
var ErrNoOptionsAcceptedForMessenger = common.NewErrorWithCategories(
	"no options accepted for messenger", common.ErrTypeMessengers,
)

var ErrWrapperFormatMessage = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeFormatMessage)
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).
	SetChild(common.ErrWrapperDatabase)

// Job errors
// Messengers wrap the Handle method of the job. If there's an error from the outer handler, we use this wrapper
var ErrWrapperHandlerWrapper = common.NewErrorWrapper("handler wrapper")
