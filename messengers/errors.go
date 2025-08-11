package messengers

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeSend = "send"
	// Lower level
	ErrTypeFormatMessage = "format message"
	ErrTypeReadUser      = "read user"
	ErrTypePrepare       = "prepare"
	ErrTypeEnqueueJob    = "enqueue job"
)

var ErrUnknownMessengerType = common.NewErrorWithCategories(
	common.ErrTypeMessengers, "unknown messenger type",
)
var ErrUserNotFound = common.NewErrorWithCategories(
	common.ErrTypeMessengers, "user not found",
)

var ErrWrapperReadUser = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypeReadUser,
)
var ErrWrapperPrepare = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypePrepare,
)
var ErrWrapperEnqueueJob = common.NewErrorWrapper(
	common.ErrTypeMessengers, ErrTypeEnqueueJob,
)
var ErrWrapperSend = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeSend)

var ErrWrapperFormat = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeFormatMessage)
var ErrWrapperAPI = common.NewErrorWrapper(common.ErrTypeMessengers, common.ErrTypeAPI)
var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).
	SetChild(common.ErrWrapperDatabase)
