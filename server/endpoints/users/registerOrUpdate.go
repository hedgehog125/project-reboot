package users

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type RegisterPayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Content  string `json:"content"  binding:"required,min=1,max=100000000"` // 100 MB but base64 encoded
	Filename string `json:"filename" binding:"required,min=1,max=256"`
	Mime     string `json:"mime" binding:"required,min=1,max=256"`
}
type RegisterOrUpdateResponse struct {
	Errors []string `json:"errors" binding:"required"`
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

		// TODO: delete active attempts if this is an update

		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, RegisterOrUpdateResponse{
			Errors: []string{},
		})
	}
}
