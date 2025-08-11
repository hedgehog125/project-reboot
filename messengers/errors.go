package messengers

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeSend = "send"
	// Lower level
	ErrTypeFormatMessage = "format message"
)

var ErrWrapperFormat = common.NewErrorWrapper(common.ErrTypeMessengers, ErrTypeFormatMessage)
var ErrWrapperAPI = common.NewErrorWrapper(common.ErrTypeMessengers, common.ErrTypeAPI)
