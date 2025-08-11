package messengers

import (
	"log"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
)

type Registry struct {
	messengers map[string]*Definition
	App        *common.App
}

type Definition struct {
	ID      string
	Version int
	// Returns the data the Handler needs, typically just a string
	Prepare       PrepareFunc
	Handler       jobs.HandlerFunc
	jobDefinition *jobs.Definition
}

// TODO: create ErrMessengerDisabledForUser
type PrepareFunc = func(user *ent.User) (any, error)

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		messengers: make(map[string]*Definition),
		App:        app,
	}
}

// TODO: how can this be called by the job definitions package?
func (registry *Registry) Register(definition *Definition) {
	fullID := messengerscommon.GetVersionedType(definition.ID, definition.Version)
	if _, exists := registry.messengers[fullID]; exists {
		log.Fatalf("messenger definition with ID \"%s\" already exists", fullID)
	}
	definition.jobDefinition = &jobs.Definition{} // TODO
	registry.messengers[fullID] = definition
}
func (registry *Registry) RegisterJobs(group *jobs.RegistryGroup) {
	for _, messenger := range registry.messengers {
		group.Register(messenger.jobDefinition)
	}
}
