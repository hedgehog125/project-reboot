package messengers

import (
	"github.com/hedgehog125/project-reboot/common"
)

type Messenger interface {
	Id() string
	Send(message common.Message) error
}
