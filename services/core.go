package services

import (
	"context"
	"sync"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
)

type Core struct {
	App       *common.App
	AdminCode core.AdminCode
	mu        sync.RWMutex
}

func NewCore(app *common.App) *Core {
	return &Core{
		App: app,
		// AdminCode will be initialised by the scheduler
	}
}

func (service *Core) RotateAdminCode() {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.AdminCode = core.NewAdminCode()
	service.AdminCode.Print()
}
func (service *Core) CheckAdminCode(givenCode string) bool {
	service.mu.RLock()
	defer service.mu.RUnlock()
	return core.CheckAdminCode(givenCode, service.AdminCode, service.App.Logger)
}
func (service *Core) SendActiveSessionReminders(ctx context.Context) *common.Error {
	return core.SendActiveSessionReminders(
		ctx, service.App.Clock, service.App.Messengers,
	)
}
func (service *Core) DeleteExpiredSessions(ctx context.Context) *common.Error {
	return core.DeleteExpiredSessions(ctx, service.App.Clock)
}

func (service *Core) Encrypt(data []byte, encryptionKey []byte) ([]byte, []byte, *common.Error) {
	return core.Encrypt(data, encryptionKey)
}

func (service *Core) Decrypt(encrypted []byte, encryptionKey []byte, nonce []byte) ([]byte, *common.Error) {
	return core.Decrypt(encrypted, encryptionKey, nonce)
}

func (service *Core) GenerateSalt() []byte {
	return core.GenerateSalt()
}

func (service *Core) HashPassword(password string, salt []byte, settings *common.PasswordHashSettings) []byte {
	return core.HashPassword(password, salt, settings)
}
