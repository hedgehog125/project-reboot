package definitions

import (
	"github.com/hedgehog125/project-reboot/messengers"
)

func Register(registry *messengers.Registry) {
	env := registry.App.Env

	if env.ENABLE_DEVELOP_MESSENGER {
		registry.Register(Develop1())
	}
	if env.DISCORD_TOKEN != "" {
		registry.Register(Discord1(registry.App))
	}
	if env.SENDGRID_TOKEN != "" {
		// TODO
	}
}
