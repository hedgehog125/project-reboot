package setup

import (
	"net/http"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

type GenerateConstantsPayload struct {
	Password string `binding:"required,min=8,max=256" json:"password"`
}

type GenerateConstantsResponse struct {
	Errors  []servercommon.ErrorDetail `binding:"required" json:"errors"`
	EnvVars *common.AdminAuthEnvVars   `                   json:"envVars"`
	TotpURL string                     `                   json:"totpUrl"`
}

func GenerateConstants(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := GenerateConstantsPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		envVars, totpURL, wrappedErr := app.Setup.GenerateAdminSetupConstants(body.Password)
		if wrappedErr != nil {
			return wrappedErr
		}

		ginCtx.JSON(http.StatusOK, GenerateConstantsResponse{
			Errors:  []servercommon.ErrorDetail{},
			EnvVars: envVars,
			TotpURL: totpURL,
		})
		return nil
	})
}
