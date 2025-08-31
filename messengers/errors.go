package messengers

import (
	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeSend = "send"
	// Lower level
	ErrTypeFormatMessage = "format message"
	ErrTypePrepare       = "prepare"
	ErrTypeEnqueueJob    = "enqueue job"
)

var ErrUnknownMessengerType = common.NewErrorWithCategories(
	"unknown messenger type", common.ErrTypeMessengers,
)

// Note: this is returned by messengers individually, some may have the contacts for the user
var ErrNoContactForUser = common.NewErrorWithCategories(
	"messenger type disabled for user", common.ErrTypeMessengers,
)

var ErrWrapperPrepare = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypePrepare,
)
var ErrWrapperEnqueueJob = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypeEnqueueJob,
)
var ErrWrapperSend = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeSend)

var ErrWrapperFormat = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeFormatMessage)
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).
	SetChild(common.ErrWrapperDatabase)

// Job errors
var ErrWrapperHandlerWrapper = common.NewErrorWrapper("handler wrapper")
