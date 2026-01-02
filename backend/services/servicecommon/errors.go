package servicecommon

import "github.com/NicoClack/cryptic-stash/backend/common"

// TODO: move types like common.Message to this package?

const (
	ErrTypeSendMessage = "send message"
)

var ErrWrapperSendMessage = common.NewErrorWrapper(ErrTypeSendMessage, common.ErrTypeMessengers)
