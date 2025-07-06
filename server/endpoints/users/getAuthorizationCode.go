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
	Errors                   []servercommon.ErrorDetail `binding:"required" json:"errors"`
	MessagesSent             []string                   `json:"messagesSent"`
	AuthorizationCode        string                     `json:"authorizationCode"`
	AuthorizationCodeValidAt time.Time                  `json:"authorizationCodeValidAt"`
}

func GetAuthorizationCode(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.Database.Client()
	messenger := app.Messenger
	clock := app.Clock
	unlockTime := time.Duration(app.Env.UNLOCK_TIME) * time.Second

	return func(ctx *gin.Context) {
		body := GetAuthorizationCodePayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			ctx.Error(ctxErr)
			return
		}

		userRow, stdErr := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Only(context.Background())
		if stdErr != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(stdErr))
			return
		}

		encryptionKey := core.HashPassword(
			body.Password,
			userRow.KeySalt,
			&common.PasswordHashSettings{
				Time:    userRow.HashTime,
				Memory:  userRow.HashMemory,
				Threads: userRow.HashThreads,
			},
		)
		_, commErr := core.Decrypt(encryptionKey, userRow.Content, userRow.Nonce)
		if commErr != nil {
			ctx.Error(servercommon.NewUnauthorizedError())
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

		authCode := core.RandomAuthCode()
		validAt := clock.Now().Add(unlockTime)

		_, stdErr = dbClient.Session.Create().
			SetUser(userRow).
			SetCode(authCode).
			SetCodeValidFrom(validAt).
			SetUserAgent(ctx.Request.UserAgent()).
			SetIP(ctx.ClientIP()).
			Save(context.Background())
		if stdErr != nil {
			ctx.Error(stdErr)
			return
		}

		ctx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
			Errors:                   []servercommon.ErrorDetail{},
			MessagesSent:             common.GetSuccessfulActionIDs(messengerIDs, errs),
			AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
			AuthorizationCodeValidAt: validAt,
		})
	}
}
