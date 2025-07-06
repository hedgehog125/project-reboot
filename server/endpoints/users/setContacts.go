package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
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
	Errors       []servercommon.ErrorDetail `binding:"required" json:"errors"`
	MessagesSent []string                   `json:"messagesSent"`
}

func SetContacts(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.Database.Client()
	messenger := app.Messenger

	return func(ctx *gin.Context) {
		body := SetContactsPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			ctx.Error(ctxErr)
			return
		}

		_, stdErr := dbClient.User.Update().
			Where(user.Username(body.Username)).
			SetAlertDiscordId(body.DiscordUserId).SetAlertEmail(body.Email).Save(ctx.Request.Context())
		if stdErr != nil {
			ctx.Error(servercommon.Send404IfNotFound(stdErr))
			return
		}

		userInfo, commErr := messengerscommon.ReadMessageUserInfo(body.Username, dbClient)
		if commErr != nil {
			ctx.Error(commErr)
			return
		}
		errs := messenger.SendUsingAll(common.Message{
			Type: common.MessageTest,
			User: userInfo,
		})
		messengerIds := messenger.IDs()
		if len(errs) == len(messengerIds) {
			ctx.JSON(http.StatusInternalServerError, SetContactsResponse{ // We aren't sure if this error is the client or server's fault
				Errors: []servercommon.ErrorDetail{
					{
						Message: "all messages failed",
						Code:    "ALL_MESSAGES_FAILED",
					},
				},
				MessagesSent: []string{},
			})
			return
		}

		// TODO: log these errors

		ctx.JSON(http.StatusOK, SetContactsResponse{
			Errors:       []servercommon.ErrorDetail{},
			MessagesSent: common.GetSuccessfulActionIDs(messengerIds, errs),
		})
	}
}
