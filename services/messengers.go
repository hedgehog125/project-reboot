package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/messengers/definitions"
)

type Messengers struct {
	*messengers.Registry
	app *common.App
}

func NewMessengers(app *common.App) *Messengers {
	registry := messengers.NewRegistry(app)
	definitions.Register(registry)

	return &Messengers{
		Registry: registry,
		app:      app,
	}
}
