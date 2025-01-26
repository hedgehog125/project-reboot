package endpoints

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
)

func RootRedirect(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	})
}

type RegisterUserPayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Content  string `json:"content"  binding:"required,min=1,max=100000000"` // 100 MB but base64 encoded
	FileName string `json:"fileName" binding:"required,min=1,max=256"`
	MimeType string `json:"mimeType" binding:"required,min=1,max=256"`
}

func RegisterUser(engine *gin.Engine, adminMiddleware gin.HandlerFunc, dbClient *ent.Client) {
	engine.POST("/api/v1/register", adminMiddleware, func(ctx *gin.Context) {
		body := RegisterUserPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		contentBytes, err := base64.RawStdEncoding.DecodeString(body.Content)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": []string{"CONTENT_DECODE_ERROR"},
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
			SetFileName(body.FileName).
			SetMime(body.MimeType).
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
			fmt.Printf("warning: an error occurred while saving user data:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": []string{"INTERNAL"},
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"errors": []string{},
		})
	})
}
