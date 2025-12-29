package users

import (
	"context"
	"net/http"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/ent/user"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

type RegisterPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string `binding:"required,min=8,max=256"                   json:"password"`
	Content  string `binding:"required,min=1,max=100000000"             json:"content"` // 100 MB but base64 encoded
	Filename string `binding:"required,min=1,max=256"                   json:"filename"`
	Mime     string `binding:"required,min=1,max=256"                   json:"mime"`
}
type RegisterOrUpdateResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

func RegisterOrUpdate(app *servercommon.ServerApp) gin.HandlerFunc {
	hashSettings := app.Env.PASSWORD_HASH_SETTINGS

	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := RegisterPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}
		contentBytes, ctxErr := servercommon.DecodeBase64(body.Content)
		if ctxErr != nil {
			return ctxErr
		}

		salt := app.Core.GenerateSalt()
		encryptionKey := app.Core.HashPassword(body.Password, salt, hashSettings)
		encrypted, nonce, wrappedErr := app.Core.Encrypt(contentBytes, encryptionKey)
		if wrappedErr != nil {
			return wrappedErr
		}

		return dbcommon.WithWriteTx(
			ginCtx.Request.Context(), app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				stdErr := tx.User.Create().
					SetUsername(body.Username).
					SetSessionsValidFrom(app.Clock.Now()).
					SetContent(encrypted).
					SetFileName(body.Filename).
					SetMime(body.Mime).
					SetNonce(nonce).
					SetKeySalt(salt).
					SetHashTime(hashSettings.Time).
					SetHashMemory(hashSettings.Memory).
					SetHashThreads(hashSettings.Threads).
					OnConflict().UpdateNewValues().
					Exec(ctx)
				if stdErr != nil {
					return stdErr
				}

				userOb, stdErr := tx.User.Query().
					Where(user.Username(body.Username)).
					Only(ctx)
				if stdErr != nil {
					return stdErr
				}

				wrappedErr := app.Core.InvalidateUserSessions(userOb.ID, ctx)
				if wrappedErr != nil {
					return wrappedErr
				}

				_, _, wrappedErr = app.Messengers.SendUsingAll(
					&common.Message{
						Type: common.MessageUserUpdate,
						User: userOb,
					},
					ctx,
				)
				if wrappedErr != nil {
					return wrappedErr
				}

				ginCtx.JSON(http.StatusCreated, RegisterOrUpdateResponse{
					Errors: []servercommon.ErrorDetail{},
				})
				return nil
			},
		)
	})
}
