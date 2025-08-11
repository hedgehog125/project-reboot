package definitions

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers"
)

func Develop1() *messengers.Definition {
	return &messengers.Definition{
		ID:      "develop",
		Version: 1,
		Prepare: func(message *common.Message) (any, error) {
			formattedMessage, commErr := messengers.FormatDefaultMessage(message)
			if commErr != nil {
				return nil, commErr
			}
			return fmt.Sprintf(
				"\nmessage sent to user \"%v\":\n%v\n\n",
				message.User.Username, formattedMessage,
			), nil
		},
		BodyType: "",
		Handler: func(jobCtx *jobs.Context) error {
			body := ""
			jobErr := jobCtx.Decode(&body)
			if jobErr != nil {
				return jobErr
			}

			fmt.Printf(body)
			return nil
		},
	}
}
