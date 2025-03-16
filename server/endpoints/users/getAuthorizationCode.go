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
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type GetAuthorizationCodeResponse struct {
	Errors                   []string  `json:"errors" binding:"required"`
	AuthorizationCode        string    `json:"authorizationCode"`
	AuthorizationCodeValidAt time.Time `json:"authorizationCodeValidAt"`
}

func GetAuthorizationCode(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()
	clock := app.App.Clock
	unlockTime := time.Duration(app.App.Env.UNLOCK_TIME) * time.Second

	return func(ctx *gin.Context) {
		body := GetAuthorizationCodePayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		userRow, err := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Select(user.FieldPasswordHash, user.FieldPasswordSalt, user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen).
			Only(context.Background())
		if err != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(err))
			return
		}

		if !core.CheckPassword(
			body.Password,
			userRow.PasswordHash,
			userRow.PasswordSalt,
			&core.HashSettings{
				Time:   userRow.HashTime,
				Memory: userRow.HashMemory,
				KeyLen: userRow.HashKeyLen,
			},
		) {
			ctx.Error(servercommon.NewUnauthorizedError())
			return
		}

		authCode := common.CryptoRandomBytes(core.AUTH_CODE_BYTE_LENGTH)
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
			AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
			AuthorizationCodeValidAt: validAt,
		})
	}
}
