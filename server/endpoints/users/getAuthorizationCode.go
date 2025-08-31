package users

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
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

	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := GetAuthorizationCodePayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		userOb, stdErr := dbcommon.WithReadTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
			userOb, stdErr := tx.User.Query().
				Where(user.Username(body.Username)).
				Only(ctx)
			if stdErr != nil {
				return nil, servercommon.SendUnauthorizedIfNotFound(
					common.ErrWrapperDatabase.Wrap(stdErr),
				)
			}
			return userOb, nil
		})
		if stdErr != nil {
			return stdErr
		}

		encryptionKey := core.HashPassword(
			body.Password,
			userOb.KeySalt,
			&common.PasswordHashSettings{
				Time:    userOb.HashTime,
				Memory:  userOb.HashMemory,
				Threads: userOb.HashThreads,
			},
		)
		_, commErr := core.Decrypt(userOb.Content, encryptionKey, userOb.Nonce)
		if commErr != nil {
			return servercommon.NewUnauthorizedError()
		}

		return dbcommon.WithWriteTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) error {
			_, _, commErr := app.Messengers.SendUsingAll(
				&common.Message{
					Type: common.MessageLogin,
					User: userOb,
				},
				ctx,
			)
			if commErr != nil {
				return commErr
			}

			authCode := core.RandomAuthCode()
			validAt := clock.Now().Add(app.Env.UNLOCK_TIME)

			_, stdErr = tx.Session.Create().
				SetUser(userOb).
				SetCode(authCode).
				SetCodeValidFrom(validAt).
				SetUserAgent(ginCtx.Request.UserAgent()).
				SetIP(ginCtx.ClientIP()).
				Save(ctx)
			if stdErr != nil {
				return common.ErrWrapperDatabase.Wrap(stdErr)
			}

			ginCtx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
				Errors:                   []servercommon.ErrorDetail{},
				AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
				AuthorizationCodeValidAt: validAt,
			})
			return nil
		})
	})
}
