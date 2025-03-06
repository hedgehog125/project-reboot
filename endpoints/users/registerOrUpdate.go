package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
)

type RegisterPayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Content  string `json:"content"  binding:"required,min=1,max=100000000"` // 100 MB but base64 encoded
	Filename string `json:"filename" binding:"required,min=1,max=256"`
	Mime     string `json:"mime" binding:"required,min=1,max=256"`
}

func RegisterOrUpdate(engine *gin.Engine, dbClient *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := RegisterPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		contentBytes, err := base64.StdEncoding.DecodeString(body.Content)
		if err != nil {
			fmt.Printf("err.Error(): %v\n", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"MALFORMED_CONTENT"},
			})
			return
		}

		encrypted, err := core.Encrypt(contentBytes, body.Password)
		if err != nil {
			fmt.Printf("warning: an error occurred while encrypting a user's data:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": []string{"INTERNAL"},
			})
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
			fmt.Printf("warning: an error occurred while saving user data:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": []string{"INTERNAL"},
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"errors": []string{},
		})
	}
}
