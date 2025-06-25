package twofactoractions

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

type Registry struct {
	actions map[string]ActionDefinition
	App     *common.App
}

type HandlerFunc func(action *Action) *common.Error
type ActionDefinition struct {
	ID       string
	Version  int
	Handler  HandlerFunc
	BodyType func() any
}

type Action struct {
	Definition *ActionDefinition
	ExpiresAt  time.Time
	Context    context.Context
	Body       []byte
}

func (action *Action) Decode(pointer any) *common.Error {
	err := json.Unmarshal(action.Body, pointer)
	if err != nil {
		return ErrWrapperDecodeAction.Wrap(err)
	}
	return nil
}

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		actions: make(map[string]ActionDefinition),
		App:     app,
	}
}

func (registry *Registry) RegisterAction(action ActionDefinition) {
	fullID := GetVersionedType(action.ID, action.Version)
	if _, exists := registry.actions[fullID]; exists {
		log.Fatalf("action with ID \"%s\" already exists", fullID)
	}
	registry.actions[fullID] = action
}
