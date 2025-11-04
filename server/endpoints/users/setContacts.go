package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type SetContactsPayload struct {
	Username      string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	DiscordUserId string `binding:"max=256"                                  json:"discordUserId"`
	Email         string `binding:"max=256"                                  json:"email"`
}
type SetContactsResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

func SetContacts(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := SetContactsPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		return dbcommon.WithWriteTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) error {
			userOb, stdErr := tx.User.Query().
				Where(user.Username(body.Username)).
				Only(ctx)
			if stdErr != nil {
				return servercommon.Send404IfNotFound(stdErr)
			}
			userOb, stdErr = userOb.Update().
				SetAlertDiscordId(body.DiscordUserId).
				SetAlertEmail(body.Email).
				Save(ctx)
			if stdErr != nil {
				return servercommon.Send404IfNotFound(stdErr)
			}

			_, _, wrappedErr := app.Messengers.SendUsingAll(
				&common.Message{
					Type: common.MessageTest,
					User: userOb,
				},
				ctx,
			)
			if wrappedErr != nil {
				return wrappedErr
			}

			ginCtx.JSON(http.StatusOK, SetContactsResponse{
				Errors: []servercommon.ErrorDetail{},
			})
			return nil
		})
	})
}
