package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/jobs/definitions"
)

type Jobs struct {
	*jobs.Engine
}

func NewJobs(app *common.App, registerFuncs ...func(group *jobs.RegistryGroup)) *Jobs {
	registry := jobs.NewRegistry(app)
	definitions.Register(registry.Group(""))
	for _, registerFunc := range registerFuncs {
		registerFunc(registry.Group(""))
	}

	return &Jobs{
		Engine: jobs.NewEngine(registry),
	}
}

func (service *Jobs) Start() {
	go service.Engine.Listen()
}

// TODO: is this the best approach?
func (service *Jobs) Encode(versionedType string, data any) (string, *common.Error) {
	return service.Engine.Registry.Encode(versionedType, data)
}
