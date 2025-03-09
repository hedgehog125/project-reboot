package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
)

func NewMessenger(env *common.Env) common.MessengerService {
	enabledMessengers := []messengers.Messenger{}
	if env.DISCORD_TOKEN != "" {
		enabledMessengers = append(enabledMessengers, messengers.NewDiscord(env))
	}
	if env.SENDGRID_TOKEN != "" {
		// TODO
	}

	return &messengerService{
		messengers: enabledMessengers,
	}
}

type messengerService struct {
	messengers []messengers.Messenger
}

func (service *messengerService) SendUsingAll(message common.Message) []error {
	errChan := make(chan common.ErrWithIndex, 3)
	for i, messenger := range service.messengers {
		go func() {
			err := messenger.Send(message)
			errChan <- common.ErrWithIndex{
				Err:   err,
				Index: i,
			}
		}()
	}
	errs := make([]error, len(service.messengers))
	for range len(service.messengers) {
		errInfo := <-errChan
		errs[errInfo.Index] = errInfo.Err
	}

	return errs
}
