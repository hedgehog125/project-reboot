package messengers

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
)

type develop struct{}

type Develop interface {
	Messenger
}

func NewDevelop() Develop {
	return &develop{}
}

func (develop *develop) Id() string {
	return "develop"
}

func (develop *develop) Send(message common.Message) *common.Error {
	formattedMessage, err := formatDefaultMessage(message)
	if err != nil {
		return err.AddCategory(ErrTypeSend)
	}

	fmt.Printf("\nmessage sent to user \"%v\":\n%v\n\n", message.User.Username, formattedMessage)

	return nil
}
