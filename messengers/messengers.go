package messengers

import "github.com/hedgehog125/project-reboot/common"

type Messenger interface {
	Send(message common.Message) error
}
