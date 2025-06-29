package users

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
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
	dbClient := app.App.Database.Client()

	return func(ctx *gin.Context) {
		body := RegisterPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			ctx.Error(ctxErr)
			return
		}

		contentBytes, stdErr := base64.StdEncoding.DecodeString(body.Content)
		if stdErr != nil {
			ctx.JSON(http.StatusBadRequest, RegisterOrUpdateResponse{
				Errors: []servercommon.ErrorDetail{
					{
						Message: "content is not valid base64",
						Code:    "MALFORMED_CONTENT",
					},
				},
			})
			return
		}

		encrypted, commErr := core.Encrypt(contentBytes, body.Password)
		if commErr != nil {
			ctx.Error(commErr)
			return
		}

		stdErr = dbClient.User.Create().
			SetUsername(body.Username).
			SetContent(encrypted.Data).
			SetFileName(body.Filename).
			SetMime(body.Mime).
			SetNonce(encrypted.Nonce).
			SetKeySalt(encrypted.KeySalt).
			SetPasswordHash(encrypted.PasswordHash).
			SetPasswordSalt(encrypted.PasswordSalt).
			SetHashTime(encrypted.HashSettings.Time).
			SetHashMemory(encrypted.HashSettings.Memory).
			SetHashKeyLen(encrypted.HashSettings.KeyLen).
			OnConflict().UpdateNewValues().
			Exec(context.Background())
		if stdErr != nil {
			ctx.Error(stdErr)
			return
		}

		userInfo, commErr := messengerscommon.ReadMessageUserInfo(body.Username, dbClient)
		if commErr == nil {
			// TODO: if this fails, let the user know using other methods
			_ = app.App.Messenger.SendUsingAll(common.Message{
				Type: common.MessageReset,
				User: userInfo,
			})
		}

		_, stdErr = dbClient.Session.Delete().Where(
			session.HasUserWith(user.Username(body.Username)),
		).Exec(context.Background())
		if stdErr != nil {
			ctx.Error(stdErr)
			return
		}

		ctx.JSON(http.StatusCreated, RegisterOrUpdateResponse{
			Errors: []servercommon.ErrorDetail{},
		})
	}
}
