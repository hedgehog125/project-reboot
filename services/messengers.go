package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
)

func NewMessenger(env *common.Env) common.MessengerService {
	enabledMessengers := []messengers.Messenger{}
	if env.IS_DEV {
		enabledMessengers = append(enabledMessengers, messengers.NewDevelop())
	}
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

func (service *messengerService) IDs() []string {
	ids := make([]string, len(service.messengers))
	for i, messenger := range service.messengers {
		ids[i] = messenger.Id()
	}
	return ids
}

// TODO: how can this be adapted to work better with common.Error?
func (service *messengerService) SendUsingAll(message common.Message) []*common.ErrWithStrId {
	errChan := make(chan common.ErrWithStrId, 3)
	for _, messenger := range service.messengers {
		go func() {
			commErr := messenger.Send(message)
			errChan <- common.ErrWithStrId{
				Err: commErr.StandardError(),
				Id:  messenger.Id(),
			}
		}()
	}
	errs := make([]*common.ErrWithStrId, 0)
	for range len(service.messengers) {
		errInfo := <-errChan
		if errInfo.Err != nil {
			errs = append(errs, &errInfo)
		}
	}

	return errs
}
