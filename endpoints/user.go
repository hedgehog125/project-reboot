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
	"github.com/hedgehog125/project-reboot/ent/loginattempt"
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

			var authCode []byte = nil
			if authorized {
				authCode = util.CryptoRandomBytes(core.AUTH_CODE_BYTE_LENGTH)
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
				"authorizationCode":        base64.StdEncoding.EncodeToString(authCode),
				"authorizationCodeValidAt": validAt,
				"content":                  nil,
				"filename":                 nil,
				"mime":                     nil,
			})
		} else {
			givenAuthCodeBytes, err := base64.StdEncoding.DecodeString(body.AuthorizationCode)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"errors": []string{"MALFORMED_AUTH_CODE"},
				})
				return
			}

			attemptRow, err := dbClient.LoginAttempt.Query().
				Where(loginattempt.And(loginattempt.Username(body.Username), loginattempt.Code(givenAuthCodeBytes))).
				Select(loginattempt.FieldCode, loginattempt.FieldCodeValidFrom).
				First(context.Background())
			if err != nil {
				if ent.IsNotFound(err) {
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"errors": []string{"INVALID_USERNAME_OR_AUTH_CODE"},
					})
				} else {
					fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"errors": []string{"INTERNAL"},
					})
				}

				return
			}

			// body.AuthorizationCode != "" so this should be a successful login attempt, but just in case
			if len(attemptRow.Code) != core.AUTH_CODE_BYTE_LENGTH {
				fmt.Printf("warning: row.Code was the wrong length! this shouldn't happen. len(row.Code): %v\n", len(attemptRow.Code))
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{"INTERNAL"},
				})
				return
			}
			if clock.Now().UTC().Before(attemptRow.CodeValidFrom) {
				ctx.JSON(http.StatusConflict, gin.H{
					"errors":                   []string{"CODE_NOT_VALID_YET"},
					"authorizationCodeValidAt": attemptRow.CodeValidFrom,
				})
				return
			}

			userRow, err := dbClient.User.Query().
				Where(user.Username(body.Username)).
				Select(user.Columns...).
				Only(context.Background())
			if err != nil {
				fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{"INTERNAL"},
				})
				return
			}

			decrypted, err := core.Decrypt(body.Password, &core.EncryptedData{
				Data:         userRow.Content,
				Nonce:        userRow.Nonce,
				KeySalt:      userRow.KeySalt,
				PasswordHash: userRow.PasswordHash,
				PasswordSalt: userRow.PasswordSalt,
				HashSettings: core.HashSettings{
					Time:   userRow.HashTime,
					Memory: userRow.HashMemory,
					KeyLen: userRow.HashKeyLen,
				},
			})
			if err != nil {
				if err == core.ErrIncorrectPassword {
					ctx.JSON(http.StatusUnauthorized, gin.H{
						"errors": []string{"INVALID_PASSWORD"},
					})
				} else {
					fmt.Printf("warning: an error occurred while decrypting user data:\n%v\n", err.Error())
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"errors": []string{"INTERNAL"},
					})
				}
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"errors":                   []string{},
				"authorizationCode":        nil,
				"authorizationCodeValidAt": nil,
				"content":                  decrypted,
				"filename":                 userRow.FileName,
				"mime":                     userRow.Mime,
			})

			// TODO: log this event to database
		}
	})
}

type RegisterUserPayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Content  string `json:"content"  binding:"required,min=1,max=100000000"` // 100 MB but base64 encoded
	Filename string `json:"filename" binding:"required,min=1,max=256"`
	Mime     string `json:"mime" binding:"required,min=1,max=256"`
}

func RegisterUser(engine *gin.Engine, adminMiddleware gin.HandlerFunc, dbClient *ent.Client) {
	engine.POST("/api/v1/users/register-or-update", adminMiddleware, func(ctx *gin.Context) {
		body := RegisterUserPayload{}
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
	})
}

type SetUserContactsPayload struct {
	Username      string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	DiscordUserId string `json:"discordUserId" binding:"max=256"`
	Email         string `json:"email" binding:"max=256"`
}

func SetUserContacts(engine *gin.Engine, adminMiddleware gin.HandlerFunc, dbClient *ent.Client) {
	engine.POST("/api/v1/users/set-user-contacts", adminMiddleware, func(ctx *gin.Context) {
		body := SetUserContactsPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		_, err := dbClient.User.Update().
			Where(user.Username(body.Username)).
			SetAlertDiscordId(body.DiscordUserId).SetAlertEmail(body.Email).Save(ctx.Request.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				ctx.JSON(http.StatusNotFound, gin.H{
					"errors": []string{"NO_USER"},
				})
			} else {
				fmt.Printf("warning: an error occurred while updating a user:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"errors": []string{"INTERNAL"},
				})
			}

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"errors": []string{},
		})
	})
}
