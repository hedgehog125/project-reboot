package users

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type GetAuthorizationCodePayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string `binding:"required,min=8,max=256"                   json:"password"`
}

type GetAuthorizationCodeResponse struct {
	Errors            []servercommon.ErrorDetail `binding:"required" json:"errors"`
	AuthorizationCode string                     `json:"authorizationCode"`
	ValidFrom         time.Time                  `json:"validFrom"`
	ValidUntil        time.Time                  `json:"validUntil"`
}

func GetAuthorizationCode(app *servercommon.ServerApp) gin.HandlerFunc {
	clock := app.Clock

	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := GetAuthorizationCodePayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		userOb, stdErr := dbcommon.WithReadTx(
			ginCtx, app.Database,
			func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
				userOb, stdErr := tx.User.Query().
					Where(user.Username(body.Username)).
					Only(ctx)
				if stdErr != nil {
					return nil, servercommon.SendUnauthorizedIfNotFound(stdErr)
				}
				return userOb, nil
			},
		)
		if stdErr != nil {
			return stdErr
		}

		encryptionKey := app.Core.HashPassword(
			body.Password,
			userOb.KeySalt,
			&common.PasswordHashSettings{
				Time:    userOb.HashTime,
				Memory:  userOb.HashMemory,
				Threads: userOb.HashThreads,
			},
		)
		time.Sleep(10 * time.Second)
		_, wrappedErr := app.Core.Decrypt(userOb.Content, encryptionKey, userOb.Nonce)
		if wrappedErr != nil {
			return servercommon.NewUnauthorizedError()
		}

		return dbcommon.WithWriteTx(
			ginCtx, app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				authCode := app.Core.RandomAuthCode()
				validFrom := clock.Now().Add(app.Env.UNLOCK_TIME)
				validUntil := clock.Now().Add(app.Env.AUTH_CODE_VALID_FOR)

				sessionOb, stdErr := tx.Session.Create().
					SetUser(userOb).
					SetCode(authCode).
					SetValidFrom(validFrom).
					SetValidUntil(validUntil).
					SetUserAgent(ginCtx.Request.UserAgent()).
					SetIP(ginCtx.ClientIP()).
					Save(ctx)
				if stdErr != nil {
					return stdErr
				}

				_, _, wrappedErr := app.Messengers.SendUsingAll(
					&common.Message{
						Type:       common.MessageLogin,
						User:       userOb,
						Time:       validFrom,
						SessionIDs: []int{sessionOb.ID},
					},
					ctx,
				)
				if wrappedErr != nil {
					return wrappedErr
				}

				ginCtx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
					Errors:            []servercommon.ErrorDetail{},
					AuthorizationCode: base64.StdEncoding.EncodeToString(authCode),
					ValidFrom:         validFrom,
					ValidUntil:        validUntil,
				})
				return nil
			},
		)
	})
}
