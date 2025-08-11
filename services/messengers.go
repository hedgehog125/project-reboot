package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/services/servicecommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type Messengers struct {
	app        *common.App
	messengers []messengers.Messenger
}

// TODO: could work similarly to the twofactoractions service but could still have its own definitions so SendUsingAll knows what jobs to call?
// Each definition could have a callback to determine if the user has this messenger enabled?
// Registering a messenger should also register its own job? The Handler method gets called directly?
// Can I remove ReadMessageInfo somehow?

func NewMessenger(app *common.App) *Messengers {
	enabledMessengers := []messengers.Messenger{}
	if app.Env.IS_DEV {
		enabledMessengers = append(enabledMessengers, messengers.NewDevelop())
	}
	if app.Env.DISCORD_TOKEN != "" {
		enabledMessengers = append(enabledMessengers, messengers.NewDiscord(app.Env))
	}
	if app.Env.SENDGRID_TOKEN != "" {
		// TODO
	}

	return &Messengers{
		app:        app,
		messengers: enabledMessengers,
	}
}

func (service *Messengers) IDs() []string {
	ids := make([]string, len(service.messengers))
	for i, messenger := range service.messengers {
		ids[i] = messenger.Id()
	}
	return ids
}

func (service *Messengers) SendUsingAll(message common.Message, ctx context.Context) *common.Error {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return servicecommon.ErrWrapperSendMessage.Wrap(common.ErrNoTxInContext)
	}

	// TODO: change messengers from interfaces to a definition based system
	jobID, commErr := service.app.Jobs.Enqueue(
		"users/SEND_MESSAGE_1",
		action.Data,
		ctx,
	)
	if commErr != nil {
		return uuid.UUID{}, twofactoractions.ErrWrapperConfirm.Wrap(commErr)
	}

	return nil
}

// Not part of service interface

func (service *Messengers) RegisterJobs(group *jobs.RegistryGroup) {
	// TODO
}
