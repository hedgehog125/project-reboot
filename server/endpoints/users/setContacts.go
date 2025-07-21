package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
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
	return servercommon.WithTx(app, func(ctx *gin.Context, tx *ent.Tx) error {
		body := SetContactsPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			return ctxErr
		}

		_, stdErr := tx.User.Update().
			Where(user.Username(body.Username)).
			SetAlertDiscordId(body.DiscordUserId).SetAlertEmail(body.Email).Save(ctx)
		if stdErr != nil {
			return servercommon.Send404IfNotFound(stdErr)
		}

		userInfo, commErr := messengerscommon.ReadUserContacts(body.Username, ctx)
		if commErr != nil {
			return commErr
		}
		commErr = app.Messengers.SendUsingAll(common.Message{
			Type: common.MessageTest,
			User: userInfo,
		})
		if commErr != nil {
			return commErr
		}

		// TODO: log these errors

		ctx.JSON(http.StatusOK, SetContactsResponse{
			Errors: []servercommon.ErrorDetail{},
		})
		return nil
	})
}
