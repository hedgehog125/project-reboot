package messengers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/jobs/jobscommon"
)

const JobNamePrefix = "messengers"

type Registry struct {
	App        *common.App
	messengers map[string]*Definition
}

type Definition struct {
	ID      string
	Version int
	// Returns the data the Handler needs, typically a struct containing the formatted message and some sort of contact (e.g a username)
	// If the user doesn't have the right contacts for this messenger, return messengers.ErrNoContactForUser.Clone()
	Prepare PrepareFunc
	// The return type of Prepare
	BodyType      any
	Handler       jobs.HandlerFunc
	jobDefinition *jobs.Definition
}

type PrepareFunc = func(message *common.Message) (any, error)

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
	prefixedGroup := group.Group(JobNamePrefix)
	for _, messenger := range registry.messengers {
		prefixedGroup.Register(messenger.jobDefinition)
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

	preparedData, stdErr := messengerDef.Prepare(message)
	if stdErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperPrepare.Wrap(stdErr))
	}

	_, commErr := registry.App.Jobs.Enqueue(
		jobscommon.JoinPaths(JobNamePrefix, versionedType),
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
) (map[string]*common.Error, *common.Error) {
	errs := make(map[string]*common.Error)
	for versionedType := range registry.messengers {
		commErr := registry.Send(versionedType, message, ctx)
		if commErr != nil {
			errs[versionedType] = commErr
			// TODO: this method doesn't work <======
			if !ErrWrapperPrepare.HasWrapped(commErr) {
				return errs, commErr
			}
			if !errors.Is(commErr, ErrNoContactForUser) {
				// Just log an error and let the admin deal with this, there's not much the user can do
				fmt.Printf("failed to prepare message for %s: %v", versionedType, commErr)
			}
		}
	}
	return errs, nil
}
