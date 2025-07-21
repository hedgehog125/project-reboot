package users

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type GetAuthorizationCodePayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string `binding:"required,min=8,max=256"                   json:"password"`
}

type GetAuthorizationCodeResponse struct {
	Errors                   []servercommon.ErrorDetail `binding:"required" json:"errors"`
	AuthorizationCode        string                     `json:"authorizationCode"`
	AuthorizationCodeValidAt time.Time                  `json:"authorizationCodeValidAt"`
}

func GetAuthorizationCode(app *servercommon.ServerApp) gin.HandlerFunc {
	clock := app.Clock

	return servercommon.WithTx(app, func(ctx *gin.Context, tx *ent.Tx) error {
		body := GetAuthorizationCodePayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			return ctxErr
		}

		userRow, stdErr := tx.User.Query().
			Where(user.Username(body.Username)).
			Only(ctx)
		if stdErr != nil {
			return servercommon.SendUnauthorizedIfNotFound(stdErr)
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
			return servercommon.NewUnauthorizedError()
		}

		commErr = app.Messengers.SendUsingAll(common.Message{
			Type: common.MessageLogin,
			User: (
			//exhaustruct:enforce
			&common.UserContacts{
				Username:       body.Username,
				AlertDiscordId: userRow.AlertDiscordId,
				AlertEmail:     userRow.AlertEmail,
			}),
		})
		if commErr != nil {
			return commErr
		}

		authCode := core.RandomAuthCode()
		validAt := clock.Now().Add(app.Env.UNLOCK_TIME)

		_, stdErr = tx.Session.Create().
			SetUser(userRow).
			SetCode(authCode).
			SetCodeValidFrom(validAt).
			SetUserAgent(ctx.Request.UserAgent()).
			SetIP(ctx.ClientIP()).
			Save(ctx)
		if stdErr != nil {
			return stdErr
		}

		ctx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
			Errors:                   []servercommon.ErrorDetail{},
			AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
			AuthorizationCodeValidAt: validAt,
		})
		return nil
	})
}
