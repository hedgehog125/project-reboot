package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type SetContactsPayload struct {
	Username      string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	DiscordUserId string `json:"discordUserId" binding:"max=256"`
	Email         string `json:"email" binding:"max=256"`
}
type SetContactsResponse struct {
	Errors       []string `json:"errors" binding:"required"`
	MessagesSent []string `json:"messagesSent"`
}

func SetContacts(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()
	messenger := app.App.Messenger

	return func(ctx *gin.Context) {
		body := SetContactsPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		_, err := dbClient.User.Update().
			Where(user.Username(body.Username)).
			SetAlertDiscordId(body.DiscordUserId).SetAlertEmail(body.Email).Save(ctx.Request.Context())
		if err != nil {
			ctx.Error(servercommon.Send404IfNotFound(err))
			return
		}

		// TODO: wrap these errors with context
		userInfo, err := common.ReadMessageUserInfo(body.Username, dbClient)
		if err != nil {
			ctx.Error(err)
			return
		}
		errs := messenger.SendUsingAll(common.Message{
			Type: common.MessageTest,
			User: userInfo,
		})
		if common.HasErrors(errs) {
			ctx.JSON(http.StatusBadRequest, SetContactsResponse{
				Errors:       []string{"SOME_TEST_MESSAGE_FAILED"}, // TODO: say which ones
				MessagesSent: []string{},                           // TODO
			})
			return
		}

		ctx.JSON(http.StatusOK, SetContactsResponse{
			Errors:       []string{},
			MessagesSent: []string{}, // TODO
		})
	}
}
