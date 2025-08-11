package jobs

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hedgehog125/project-reboot/common"
)

type Registry struct {
	jobs map[string]*Definition
	App  *common.App
}

const (
	LowPriority      = -1
	DefaultPriority  = 0
	HighPriority     = 1
	CriticalPriority = 2
	// TODO: implement
	RealtimePriority = 3 // No limit on number of RealtimePriority jobs that can run concurrently
)

type Definition struct {
	ID                string
	Version           int
	Handler           HandlerFunc
	BodyType          any
	reflectedBodyType reflect.Type
	Weight            int
	// 0 is DefaultPriority.
	Priority int8
}
type HandlerFunc func(jobCtx *Context) error

type Context struct {
	Definition *Definition
	Context    context.Context
	Body       []byte
}

func (ctx *Context) Decode(pointer any) *Error {
	err := json.Unmarshal(ctx.Body, pointer)
	if err != nil {
		return NewError(ErrWrapperDecode.Wrap(err))
	}
	return nil
}

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		jobs: make(map[string]*Definition),
		App:  app,
	}
}

func (registry *Registry) Register(definition *Definition) {
	fullID := common.GetVersionedType(definition.ID, definition.Version)
	if _, exists := registry.jobs[fullID]; exists {
		log.Fatalf("job definition with ID \"%s\" already exists", fullID)
	}
	prepareJobDefinition(definition)
	registry.jobs[fullID] = definition
}
func prepareJobDefinition(definition *Definition) {
	fullID := common.GetVersionedType(definition.ID, definition.Version)
	if definition.BodyType != nil {
		definition.reflectedBodyType = reflect.TypeOf(definition.BodyType)
		if definition.reflectedBodyType.Kind() != reflect.Ptr {
			log.Fatalf("job definition %s body type must be a pointer", fullID)
		}
	}
	if definition.Weight < 1 {
		log.Fatalf("job definition %s weight must be 1 or higher", fullID)
	}
	if definition.Priority < LowPriority || definition.Priority > RealtimePriority {
		log.Fatalf("job definition %s priority must be between -1 (LowPriority) and 5 (RealtimePriority)", fullID)
	}
}
