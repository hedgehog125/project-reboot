package users

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
	"github.com/hedgehog125/project-reboot/server/servercommon"
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

		salt := core.GenerateSalt()
		encryptionKey := core.HashPassword(body.Password, salt, hashSettings)
		encrypted, nonce, commErr := core.Encrypt(contentBytes, encryptionKey)
		if commErr != nil {
			return commErr
		}

		return dbcommon.WithWriteTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) error {
			stdErr := tx.User.Create().
				SetUsername(body.Username).
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
				return common.ErrWrapperDatabase.Wrap(stdErr)
			}

			userInfo, commErr := messengerscommon.ReadUserContacts(body.Username, ctx)
			if commErr != nil {
				return commErr
			}
			commErr = app.Messengers.SendUsingAll(
				common.Message{
					Type: common.MessageReset,
					User: userInfo,
				},
				ctx,
			)
			if commErr != nil {
				return commErr
			}

			_, stdErr = tx.Session.Delete().Where(
				session.HasUserWith(user.Username(body.Username)),
			).Exec(ctx)
			if stdErr != nil {
				return common.ErrWrapperDatabase.Wrap(stdErr)
			}

			ginCtx.JSON(http.StatusCreated, RegisterOrUpdateResponse{
				Errors: []servercommon.ErrorDetail{},
			})
			return nil
		})
	})
}
