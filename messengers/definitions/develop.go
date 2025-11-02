package definitions

import (
	"fmt"
	"os"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
)

type Develop1Body struct {
	FullMessage string `json:"formattedMessage"`
}

func Develop1() *messengers.Definition {
	return &messengers.Definition{
		ID:      "develop",
		Version: 1,
		Prepare: func(message *common.Message) (any, error) {
			formattedMessage, commErr := messengers.FormatDefaultMessage(message)
			if commErr != nil {
				return nil, commErr
			}
			return &Develop1Body{
				FullMessage: fmt.Sprintf(
					"\nmessage sent to user \"%v\":\n%v\n",
					message.User.Username, formattedMessage,
				),
			}, nil
		},
		BodyType: &Develop1Body{},
		Handler: func(messengerCtx *messengers.Context) error {
			body := Develop1Body{}
			commErr := messengerCtx.Decode(&body)
			if commErr != nil {
				return commErr
			}

			fmt.Fprintln(os.Stdout, body.FullMessage)
			messengerCtx.ConfirmSent()
			return nil
		},
	}
}
