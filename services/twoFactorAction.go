package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/twofactoractions"
	"github.com/hedgehog125/project-reboot/twofactoractions/actions"
)

func NewTwoFactorAction(app *common.App) common.TwoFactorActionService {
	registry := twofactoractions.NewRegistry(app)
	actions.RegisterActions(registry.Group(""))

	return &jobService{
		registry: registry,
	}
}

type twoFactorActionService struct {
	registry                *twofactoractions.Registry
	runningActionsWaitGroup common.WaitGroupWithCounter
}

func (service *jobService) Shutdown() {
	select {
	case <-common.NewCallbackChannel(service.runningActionsWaitGroup.Wait):
	case <-time.After(10 * time.Second):
		fmt.Printf("warning: 2FA service timed out waiting for actions to complete during shutdown. %v are still running\n", service.runningActionsWaitGroup.WaitingCount())
	}
}

func (service *jobService) Confirm(actionID uuid.UUID, code string) *common.Error {
	service.runningActionsWaitGroup.Add(1)
	defer service.runningActionsWaitGroup.Done()

	return service.registry.Confirm(actionID, code)
}
func (service *jobService) Create(
	actionType string,
	version int,
	expiresAt time.Time,
	data any,
) (uuid.UUID, string, *common.Error) {
	return service.registry.Create(
		actionType,
		version,
		expiresAt,
		data,
	)
}
