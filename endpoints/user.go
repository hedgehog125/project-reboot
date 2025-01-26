package endpoints

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
	"github.com/jonboulle/clockwork"
)

type GetUserDownloadPayload struct {
	Username          string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password          string `json:"password" binding:"required,min=8,max=256"`
	AuthorizationCode string `json:"authorizationCode" binding:"max=256"`
}

func GetUserDownload(engine *gin.Engine, dbClient *ent.Client, clock clockwork.Clock, env *intertypes.Env) {
	engine.POST("/api/v1/users/download", func(ctx *gin.Context) {
		body := GetUserDownloadPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		if body.AuthorizationCode == "" {
			row, err := dbClient.User.Query().
				Where(user.Username(body.Username)).
				Select(user.FieldPasswordHash, user.FieldPasswordSalt, user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen).
				Only(context.Background())

			if err != nil {
				if ent.IsNotFound(err) {
					ctx.JSON(http.StatusNotFound, gin.H{
						"errors": []string{"NO_USER"},
					})
				} else {
					fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"errors": []string{"INTERNAL"},
					})
				}

				return
			}

			authorized := core.CheckPassword(
				body.Password,
				row.PasswordHash,
				row.PasswordSalt,
				&core.HashSettings{
					Time:   row.HashTime,
					Memory: row.HashMemory,
					KeyLen: row.HashKeyLen,
				},
			)

			authCode := ""
			if authorized {
				authCode = base64.RawStdEncoding.EncodeToString(
					util.CryptoRandomBytes(core.AUTH_CODE_BYTE_LENGTH),
				)
			}
			validAt := clock.Now().UTC().
				Add(time.Duration(env.UNLOCK_TIME) * time.Second)

			_, err = dbClient.LoginAttempt.Create().
				SetUsername(body.Username).
				SetCode(authCode).
				SetCodeValidFrom(validAt).
				SetInfo(&intertypes.LoginAttemptInfo{
					UserAgent: ctx.Request.UserAgent(),
					IP:        ctx.ClientIP(),
				}).Save(context.Background())
			if err != nil {
				fmt.Printf("warning: an error occurred while creating a login attempt:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{"INTERNAL"},
				})
				return
			}

			if !authorized {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"errors": []string{"INCORRECT_USERNAME_OR_PASSWORD"},
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"errors":                   []string{},
				"authorizationCode":        authCode,
				"authorizationCodeValidAt": validAt.Unix(),
				"rebootZipContent":         nil,
				"rebootZipFilename":        nil,
				"rebootZipMime":            nil,
			})
		}
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
	engine.POST("/api/v1/users/register-or-update", adminMiddleware, func(ctx *gin.Context) {
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
	})
}
