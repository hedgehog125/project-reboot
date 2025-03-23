package users

import (
	"net/http"
	"slices"

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
		userInfo, err := app.App.Database.ReadMessageUserInfo(body.Username)
		if err != nil {
			ctx.Error(err)
			return
		}
		errs := messenger.SendUsingAll(common.Message{
			Type: common.MessageTest,
			User: userInfo,
		})
		if len(errs) != 0 {
			ctx.JSON(http.StatusInternalServerError, SetContactsResponse{ // We aren't sure if this error is the client or server's fault
				Errors:       []string{"ALL_TEST_MESSAGES_FAILED"},
				MessagesSent: []string{},
			})
			return
		}

		messagesSent := messenger.Ids()
		for _, err := range errs {
			index := slices.Index(messagesSent, err.Id)
			if index != -1 {
				messagesSent = slices.Delete(messagesSent, index, index)
			}
		}

		// TODO: log these errors

		ctx.JSON(http.StatusOK, SetContactsResponse{
			Errors:       []string{},
			MessagesSent: messagesSent,
		})
	}
}
