package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/util"
)

type SetContactsPayload struct {
	Username      string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	DiscordUserId string `json:"discordUserId" binding:"max=256"`
	Email         string `json:"email" binding:"max=256"`
}

func SetContacts(engine *gin.Engine, dbClient *ent.Client, messengerGroup messengers.MessengerGroup) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := SetContactsPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		_, err := dbClient.User.Update().
			Where(user.Username(body.Username)).
			SetAlertDiscordId(body.DiscordUserId).SetAlertEmail(body.Email).Save(ctx.Request.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"errors": []string{"NO_USER"},
				})
			} else {
				fmt.Printf("warning: an error occurred while updating a user:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{"INTERNAL"},
				})
			}

			return
		}

		userInfo, err := messengers.ReadUserInfo(body.Username, dbClient)
		if err != nil {
			fmt.Printf("warning: an error occurred while reading the user info for a messenger:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": []string{"INTERNAL"},
			})
			return
		}
		errs := messengerGroup.SendUsingAll(messengers.Message{
			Type: messengers.MessageTest,
			User: userInfo,
		})
		if util.HasErrors(errs) {
			ctx.JSON(http.StatusOK, gin.H{
				"errors": []string{"TEST_MESSAGE_SEND_ERROR"}, // TODO: say which ones
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"errors": []string{},
		})
	}
}
