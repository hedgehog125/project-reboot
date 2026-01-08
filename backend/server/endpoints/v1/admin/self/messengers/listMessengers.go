package messengers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/NicoClack/cryptic-stash/backend/ent/user"
	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
	"github.com/gin-gonic/gin"
)

type ListMessengersResponse struct {
	Errors     []servercommon.ErrorDetail `binding:"required" json:"errors"`
	Messengers map[string]*Messenger      `binding:"required" json:"messengers"`
}

type Messenger struct {
	Name    string          `binding:"required" json:"name"`
	Enabled bool            `binding:"required" json:"enabled"`
	Options json.RawMessage `binding:"required" json:"options"`
}

func ListMessengers(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		configuredMessengers, stdErr := dbcommon.WithReadTx(
			ginCtx.Request.Context(), app.Database,
			func(tx *ent.Tx, ctx context.Context) (*common.SetupStatus, error) {
				userOb, stdErr := tx.User.Query().
					Where(user.Username(common.AdminUsername)).
					Only(ctx)
				if stdErr != nil {
					return nil, stdErr
				}
				return app.Messengers.GetConfiguredMessengers() // TODO: implement
			},
		)
		if stdErr != nil {
			return stdErr
		}

		ginCtx.JSON(http.StatusOK, ListMessengersResponse{
			Errors:                       []servercommon.ErrorDetail{},
			IsComplete:                   status.IsComplete,
			IsEnvComplete:                status.IsEnvComplete,
			AreAdminMessengersConfigured: status.AreAdminMessengersConfigured,
		})
		return nil
	})
}
