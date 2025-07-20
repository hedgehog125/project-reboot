package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/jobs/definitions"
)

type jobService struct {
	*jobs.Engine
}

func NewJob(app *common.App) common.JobService {
	registry := jobs.NewRegistry(app)
	definitions.Register(registry.Group(""))

	return &jobService{
		Engine: jobs.NewEngine(registry),
	}
}

func (service *jobService) Start() {
	go service.Engine.Listen()
}

// TODO: is this the best approach?
func (service *jobService) Encode(versionedType string, data any) (string, *common.Error) {
	return service.Engine.Registry.Encode(versionedType, data)
}
