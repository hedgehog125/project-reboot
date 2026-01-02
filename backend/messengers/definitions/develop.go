package definitions

import (
	"fmt"
	"os"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/messengers"
)

type Develop1Body struct {
	FullMessage string `json:"formattedMessage"`
}

func Develop1() *messengers.Definition {
	return &messengers.Definition{
		ID:      "develop",
		Version: 1,
		Prepare: func(message *common.Message) (any, error) {
			formattedMessage, wrappedErr := messengers.FormatDefaultMessage(message)
			if wrappedErr != nil {
				return nil, wrappedErr
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
			wrappedErr := messengerCtx.Decode(&body)
			if wrappedErr != nil {
				return wrappedErr
			}

			fmt.Fprintln(os.Stdout, body.FullMessage)
			messengerCtx.ConfirmSent()
			return nil
		},
	}
}
