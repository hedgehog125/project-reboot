package twofactoractions

import (
	"context"
	"log"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

type Registry struct {
	actions map[string]ActionDefinition[any]
	App     *common.App
}

type HandlerFunc[Body any] func(action *Action[Body]) *common.Error
type ActionDefinition[Body any] struct {
	ID       string
	Version  int
	Handler  HandlerFunc[Body]
	BodyType Body
}

type Action[T any] struct {
	Definition *ActionDefinition[T]
	ExpiresAt  time.Time
	Context    context.Context
	Body       T
}

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		actions: make(map[string]ActionDefinition[any]),
		App:     app,
	}
}

func (registry *Registry) RegisterAction(action ActionDefinition[any]) {
	fullID := GetVersionedType(action.ID, action.Version)
	if _, exists := registry.actions[fullID]; exists {
		log.Fatalf("action with ID \"%s\" already exists", fullID)
	}
	registry.actions[fullID] = action
}
