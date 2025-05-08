package users

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type GetAuthorizationCodePayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string `binding:"required,min=8,max=256"                   json:"password"`
}

type GetAuthorizationCodeResponse struct {
	Errors                   []string  `binding:"required"              json:"errors"`
	MessagesSent             []string  `json:"messagesSent"`
	AuthorizationCode        string    `json:"authorizationCode"`
	AuthorizationCodeValidAt time.Time `json:"authorizationCodeValidAt"`
}

func GetAuthorizationCode(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()
	messenger := app.App.Messenger
	clock := app.App.Clock
	unlockTime := time.Duration(app.App.Env.UNLOCK_TIME) * time.Second

	return func(ctx *gin.Context) {
		body := GetAuthorizationCodePayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		userRow, err := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Select(
				user.FieldPasswordHash, user.FieldPasswordSalt,
				user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen,
				// Contacts
				user.FieldAlertDiscordId,
				user.FieldAlertEmail,
			).
			Only(context.Background())
		if err != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(err))
			return
		}

		if !core.CheckPassword(
			body.Password,
			userRow.PasswordHash,
			userRow.PasswordSalt,
			core.HashSettings{
				Time:   userRow.HashTime,
				Memory: userRow.HashMemory,
				KeyLen: userRow.HashKeyLen,
			},
		) {
			ctx.Error(servercommon.ErrUnauthorized)
			return
		}

		errs := messenger.SendUsingAll(common.Message{
			Type: common.MessageLogin,
			User: (
			//exhaustruct:enforce
			&common.MessageUserInfo{
				Username:       body.Username,
				AlertDiscordId: userRow.AlertDiscordId,
				AlertEmail:     userRow.AlertEmail,
			}),
		})
		messengerIDs := messenger.IDs()
		if len(errs) == len(messengerIDs) {
			// We aren't sure if this error is the client or server's fault
			ctx.JSON(http.StatusBadRequest, SetContactsResponse{
				Errors:       []string{"ALL_TEST_MESSAGES_FAILED"},
				MessagesSent: []string{},
			})
			return
		}

		// TODO: log these errors

		authCode := core.RandomAuthCode()
		validAt := clock.Now().UTC().Add(unlockTime)

		_, err = dbClient.Session.Create().
			SetUser(userRow).
			SetCode(authCode).
			SetCodeValidFrom(validAt).
			SetUserAgent(ctx.Request.UserAgent()).
			SetIP(ctx.ClientIP()).
			Save(context.Background())
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
			Errors:                   []string{},
			MessagesSent:             common.GetSuccessfulActionIDs(messengerIDs, errs),
			AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
			AuthorizationCodeValidAt: validAt,
		})
	}
}
