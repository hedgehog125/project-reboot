package services

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/messengers/definitions"
)

const (
	BulkMessageDelay = 1500 * time.Millisecond
)

type Messengers struct {
	App      *common.App
	Registry *messengers.Registry
}

func NewMessengers(app *common.App) *Messengers {
	registry := messengers.NewRegistry(app)
	definitions.Register(registry)

	return &Messengers{
		App:      app,
		Registry: registry,
	}
}

func (service *Messengers) Send(
	versionedType string, message *common.Message,
	ctx context.Context,
) common.WrappedError {
	return service.Registry.Send(versionedType, message, service.App.Clock.Now(), ctx)
}
func (service *Messengers) ScheduleSend(
	versionedType string, message *common.Message,
	sendTime time.Time,
	ctx context.Context,
) common.WrappedError {
	return service.Registry.Send(versionedType, message, sendTime, ctx)
}
func (service *Messengers) SendUsingAll(
	message *common.Message,
	ctx context.Context,
) (int, map[string]common.WrappedError, common.WrappedError) {
	return service.Registry.SendUsingAll(message, service.App.Clock.Now(), ctx)
}
func (service *Messengers) ScheduleSendUsingAll(
	message *common.Message,
	sendTime time.Time,
	ctx context.Context,
) (int, map[string]common.WrappedError, common.WrappedError) {
	return service.Registry.SendUsingAll(message, sendTime, ctx)
}
func (service *Messengers) SendBulk(
	messages []*common.Message, ctx context.Context,
) common.WrappedError {
	return service.Registry.SendBulk(
		messages,
		func(lastSendTime time.Time, index int) time.Time {
			if lastSendTime.IsZero() {
				return service.App.Clock.Now()
			}
			return lastSendTime.Add(BulkMessageDelay)
		},
		ctx,
	)
}

// Not in service interface
func (service *Messengers) RegisterJobs(group *jobs.RegistryGroup) {
	service.Registry.RegisterJobs(group)
}
