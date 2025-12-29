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

		return dbcommon.WithWriteTx(
			ginCtx.Request.Context(), app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
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

				// The user most likely isn't trying to log in if they've coordinated this with their admin
				// And deleting the sessions simplifies IsUserSufficientlyNotified
				wrappedErr = app.Core.InvalidateUserSessions(userOb.ID, ctx)
				if wrappedErr != nil {
					return wrappedErr
				}

				ginCtx.JSON(http.StatusOK, SetContactsResponse{
					Errors: []servercommon.ErrorDetail{},
				})
				return nil
			},
		)
	})
}
