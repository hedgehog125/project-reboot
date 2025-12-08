package services

import (
	"encoding/json"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/jobs"
	"github.com/NicoClack/cryptic-stash/jobs/definitions"
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
func (service *Jobs) Encode(versionedType string, body any) (json.RawMessage, common.WrappedError) {
	return service.Engine.Registry.Encode(versionedType, body)
}
