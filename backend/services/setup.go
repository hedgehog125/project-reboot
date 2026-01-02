package services

import (
	"context"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/setup"
)

type Setup struct {
	App *common.App
}

func NewSetupService(app *common.App) *Setup {
	return &Setup{
		App: app,
	}
}

func (service *Setup) IsSetupComplete(ctx context.Context) (bool, common.WrappedError) {
	return false, nil
}

func (service *Setup) GenerateAdminSetupConstants(
	password string,
) (*common.AdminAuthEnvVars, string, common.WrappedError) {
	return setup.GenerateAdminSetupConstants(
		password,
		service.App.Env.ADMIN_PASSWORD_HASH_SETTINGS,
		service.App.Env.IS_DEV,
		service.App.Core,
	)
}

func (service *Setup) CheckTotpCode(totpCode string, totpSecret string) bool {
	return setup.CheckTotpCode(totpCode, totpSecret)
}
