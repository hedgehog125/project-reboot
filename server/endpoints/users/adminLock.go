package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type AdminLockPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
}
type AdminLockResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

// Admin route
func AdminLock(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := AdminLockPayload{}
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
					SetLocked(true).
					ClearLockedUntil().
					Save(ctx)
				if stdErr != nil {
					return servercommon.Send404IfNotFound(stdErr)
				}
				_, stdErr = tx.Session.Delete().
					Where(session.UserID(userOb.ID)).
					Exec(ctx)
				if stdErr != nil {
					return stdErr
				}

				_, _, commErr := app.Messengers.SendUsingAll(
					&common.Message{
						Type: common.MessageLock,
						User: userOb,
					},
					ctx,
				)
				if commErr != nil {
					return commErr
				}

				ginCtx.JSON(http.StatusOK, AdminLockResponse{
					Errors: []servercommon.ErrorDetail{},
				})
				return nil
			},
		)
	})
}
