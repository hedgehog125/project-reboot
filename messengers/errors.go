package messengers

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeSend = "send"
	// Lower level
	ErrTypeFormatMessage = "format message"
)

var ErrWrapperFormat = common.NewErrorWrapper(ErrTypeFormatMessage, common.ErrTypeMessengers)
var ErrWrapperAPI = common.NewErrorWrapper(common.ErrTypeAPI, common.ErrTypeMessengers)
