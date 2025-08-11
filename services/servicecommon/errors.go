package servicecommon

import "github.com/hedgehog125/project-reboot/common"

// TODO: move types like common.Message to this package?

const (
	ErrTypeSendMessage = "send message"
)

var ErrWrapperSendMessage = common.NewErrorWrapper(ErrTypeSendMessage, common.ErrTypeMessengers)
