package twofactoractions

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

type Registry struct {
	actions map[string]*ActionDefinition
	App     *common.App
}

type HandlerFunc func(action *Action) *common.Error
type ActionDefinition struct {
	ID                string
	Version           int
	Handler           HandlerFunc
	BodyType          any
	reflectedBodyType reflect.Type
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
		actions: make(map[string]*ActionDefinition),
		App:     app,
	}
}

func (registry *Registry) RegisterAction(actionDef *ActionDefinition) {
	fullID := GetVersionedType(actionDef.ID, actionDef.Version)
	if _, exists := registry.actions[fullID]; exists {
		log.Fatalf("action with ID \"%s\" already exists", fullID)
	}
	prepareActionDefinition(actionDef)
	registry.actions[fullID] = actionDef
}
func prepareActionDefinition(actionDef *ActionDefinition) {
	fullID := GetVersionedType(actionDef.ID, actionDef.Version)
	if actionDef.BodyType == nil {
		log.Fatalf("action %s has no body type", fullID)
	}
	actionDef.reflectedBodyType = reflect.TypeOf(actionDef.BodyType)
	if actionDef.reflectedBodyType.Kind() != reflect.Ptr {
		log.Fatalf("action %s body type must be a pointer", fullID)
	}
}
