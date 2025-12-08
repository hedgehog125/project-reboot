package users

import (
	"context"
	"net/http"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/ent/user"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

type AdminUnlockPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
}
type AdminUnlockResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

// Admin route
func AdminUnlock(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := AdminUnlockPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		return dbcommon.WithWriteTx(
			ginCtx, app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				userOb, stdErr := tx.User.Query().
					Where(user.Username(body.Username)).
					Only(ctx)
				if stdErr != nil {
					return servercommon.Send404IfNotFound(stdErr)
				}
				userOb, stdErr = userOb.Update().
					SetLocked(false).
					ClearLockedUntil().
					Save(ctx)
				if stdErr != nil {
					return servercommon.Send404IfNotFound(stdErr)
				}

				wrappedErr := app.Core.InvalidateUserSessions(userOb.ID, ctx)
				if wrappedErr != nil {
					return wrappedErr
				}
				_, _, wrappedErr = app.Messengers.SendUsingAll(
					&common.Message{
						Type: common.MessageUnlock,
						User: userOb,
					},
					ctx,
				)
				if wrappedErr != nil {
					return wrappedErr
				}

				ginCtx.JSON(http.StatusOK, AdminUnlockResponse{
					Errors: []servercommon.ErrorDetail{},
				})
				return nil
			},
		)
	})
}
