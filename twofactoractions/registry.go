package twofactoractions

import (
	"log"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

type Registry struct {
	actions map[string]ActionDefinition[any]
	App     *common.App
}

type ActionDefinition[Body any] struct {
	ID       string
	Version  int
	Handler  func(action *Action[Body]) *common.Error
	BodyType Body
}

type Action[T any] struct {
	Definition *ActionDefinition[T]
	ExpiresAt  time.Time
	Body       *T
}

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		actions: make(map[string]ActionDefinition[any]),
		App:     app,
	}
}

func (registry *Registry) RegisterAction(action ActionDefinition[any]) {
	fullID := GetFullType(action.ID, action.Version)
	if _, exists := registry.actions[fullID]; exists {
		log.Fatalf("action with ID \"%s\" already exists", fullID)
	}
	registry.actions[fullID] = action
}
