package messengers

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeSend = "send"
	// Lower level
	ErrTypeFormatMessage = "format message"
	ErrTypePrepare       = "prepare"
	ErrTypeEnqueueJob    = "enqueue job"
)

var ErrUnknownMessengerType = common.NewErrorWithCategories(
	common.ErrTypeMessengers, "unknown messenger type",
)

// Note: this is returned by messengers individually, some may have the contacts for the user
var ErrNoContactForUser = common.NewErrorWithCategories(
	common.ErrTypeMessengers, "messenger type disabled for user",
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
