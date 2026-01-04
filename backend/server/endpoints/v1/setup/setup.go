package setup

import (
	"context"
	"net/http"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
	"github.com/gin-gonic/gin"
)

func ConfigureEndpoints(group *servercommon.Group) {
	group.GET("/", GetSetup(group.App))
	if group.App.Env.ENABLE_ENV_SETUP {
		group.POST("/generate-constants", GenerateConstants(group.App))
		group.POST("/check-totp", CheckTotp(group.App))
		group.GET("/echo-headers", EchoHeaders(group.App))
	}
}

type GetSetupResponse struct {
	Errors                       []servercommon.ErrorDetail `binding:"required" json:"errors"`
	IsComplete                   bool                       `binding:"required" json:"isComplete"`
	IsEnvComplete                bool                       `binding:"required" json:"isEnvComplete"`
	AreAdminMessengersConfigured bool                       `binding:"required" json:"areAdminMessengersConfigured"`
}

func GetSetup(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		status, stdErr := dbcommon.WithReadTx(
			ginCtx.Request.Context(), app.Database,
			func(tx *ent.Tx, ctx context.Context) (*common.SetupStatus, error) {
				return app.Setup.GetStatus(ctx)
			},
		)
		if stdErr != nil {
			return stdErr
		}

		ginCtx.JSON(http.StatusOK, GetSetupResponse{
			Errors:                       []servercommon.ErrorDetail{},
			IsComplete:                   status.IsComplete,
			IsEnvComplete:                status.IsEnvComplete,
			AreAdminMessengersConfigured: status.AreAdminMessengersConfigured,
		})
		return nil
	})
}
