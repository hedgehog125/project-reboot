package services

import (
	"encoding/json"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/jobs"
	"github.com/NicoClack/cryptic-stash/backend/jobs/definitions"
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
	go service.Listen()
}

// TODO: is this the best approach?
func (service *Jobs) Encode(versionedType string, body any) (json.RawMessage, common.WrappedError) {
	return service.Registry.Encode(versionedType, body)
}
