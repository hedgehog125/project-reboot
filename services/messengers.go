package services

import (
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/messengers"
)

func ConfigureMessengers(env *intertypes.Env) messengers.MessengerGroup {
	enabledMessengers := []messengers.Messenger{}
	if env.DISCORD_TOKEN != "" {
		enabledMessengers = append(enabledMessengers, messengers.NewDiscord(env))
	}
	if env.SENDGRID_TOKEN != "" {
		// TODO
	}

	return messengers.NewMessengerGroup(enabledMessengers)
}
