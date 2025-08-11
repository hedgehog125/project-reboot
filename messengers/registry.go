package messengers

import (
	"context"
	"log"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
)

type Registry struct {
	App           *common.App
	messengers    map[string]*Definition
	jobNamePrefix string
}

type Definition struct {
	ID      string
	Version int
	// Returns the data the Handler needs, typically a struct containing the formatted message and some sort of contact (e.g a username)
	Prepare PrepareFunc
	// The return type of Prepare
	BodyType      any
	Handler       jobs.HandlerFunc
	jobDefinition *jobs.Definition
}

// TODO: create ErrMessengerTypeDisabledForUser
type PrepareFunc = func(message *common.Message, user *ent.User) (any, error)

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		messengers: make(map[string]*Definition),
		App:        app,
	}
}

func (registry *Registry) Register(definition *Definition) {
	fullID := common.GetVersionedType(definition.ID, definition.Version)
	if _, exists := registry.messengers[fullID]; exists {
		log.Fatalf("messenger definition with ID \"%s\" already exists", fullID)
	}
	definition.jobDefinition = &jobs.Definition{
		ID:       definition.ID,
		Version:  definition.Version,
		Handler:  definition.Handler,
		BodyType: definition.BodyType,
		Weight:   1,
	}
	registry.messengers[fullID] = definition
}
func (registry *Registry) RegisterJobs(group *jobs.RegistryGroup) {
	registry.jobNamePrefix = group.Path
	for _, messenger := range registry.messengers {
		group.Register(messenger.jobDefinition)
	}
}

func (registry *Registry) Send(
	versionedType string, message *common.Message, ctx context.Context,
) *common.Error {
	messengerDef, ok := registry.messengers[versionedType]
	if !ok {
		return ErrWrapperSend.Wrap(ErrUnknownMessengerType)
	}
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return ErrWrapperSend.Wrap(common.ErrNoTxInContext)
	}

	userRow, stdErr := tx.User.Query().Where(user.Username(message.Username)).Only(ctx)
	if stdErr != nil {
		return ErrWrapperSend.Wrap(
			ErrWrapperReadUser.Wrap(
				ErrWrapperDatabase.Wrap(stdErr),
			),
		)
	}
	preparedData, stdErr := messengerDef.Prepare(message, userRow)
	if stdErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperPrepare.Wrap(stdErr))
	}

	_, commErr := registry.App.Jobs.Enqueue(
		versionedType,
		preparedData,
		ctx,
	)
	if commErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperEnqueueJob.Wrap(commErr))
	}

	return nil
}
func (registry *Registry) SendUsingAll(
	message *common.Message, ctx context.Context,
) *common.Error {
	for versionedType, _ := range registry.messengers {
		commErr := registry.Send(versionedType, message, ctx)
		if commErr != nil {
			return commErr
		}
	}
	return nil
}
