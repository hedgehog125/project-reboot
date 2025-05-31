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
	Errors []string `binding:"required" json:"errors"`
}

func RegisterOrUpdate(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()

	return func(ctx *gin.Context) {
		body := RegisterPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		contentBytes, err := base64.StdEncoding.DecodeString(body.Content)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, RegisterOrUpdateResponse{
				Errors: []string{"MALFORMED_CONTENT"},
			})
			return
		}

		encrypted, err := core.Encrypt(contentBytes, body.Password)
		if err != nil {
			ctx.Error(err)
			return
		}

		err = dbClient.User.Create().
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
		if err != nil {
			ctx.Error(err)
			return
		}

		userInfo, err := messengerscommon.ReadMessageUserInfo(body.Username, dbClient)
		if err == nil {
			// TODO: if this fails, let the user know using other methods
			_ = app.App.Messenger.SendUsingAll(common.Message{
				Type: common.MessageReset,
				User: userInfo,
			})
		}

		_, err = dbClient.Session.Delete().Where(
			session.HasUserWith(user.Username(body.Username)),
		).Exec(context.Background())
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, RegisterOrUpdateResponse{
			Errors: []string{},
		})
	}
}
