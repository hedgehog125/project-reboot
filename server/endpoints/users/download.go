package users

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
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type DownloadPayload struct {
	Username          string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password          string `json:"password" binding:"required,min=8,max=256"`
	AuthorizationCode string `json:"authorizationCode" binding:"required,min=256,max=256"`
}

type DownloadResponse struct {
	Errors                   []string   `json:"errors" binding:"required"`
	AuthorizationCodeValidAt *time.Time `json:"authorizationCodeValidAt"`
	Content                  []byte     `json:"content"`
	Filename                 string     `json:"filename"`
	Mime                     string     `json:"mime"`
}

func Download(app *servercommon.ServerApp) gin.HandlerFunc {
	sendUnauthorizedError := func(ctx *gin.Context) {
		ctx.JSON(http.StatusUnauthorized, DownloadResponse{
			Errors: []string{"INCORRECT_USERNAME_OR_PASSWORD_OR_AUTH_CODE"},
		})
	}
	dbClient := app.App.Database.Client()
	clock := app.App.Clock

	return func(ctx *gin.Context) {
		body := DownloadPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		givenAuthCodeBytes, err := base64.StdEncoding.DecodeString(body.AuthorizationCode)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, DownloadResponse{
				Errors: []string{"MALFORMED_AUTH_CODE"},
			})
			return
		}

		attemptRow, err := dbClient.LoginAttempt.Query().
			Where(loginattempt.And(loginattempt.HasUserWith(user.Username(body.Username)), loginattempt.Code(givenAuthCodeBytes))).
			Select(loginattempt.FieldCode, loginattempt.FieldCodeValidFrom).
			First(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				sendUnauthorizedError(ctx)
			} else {
				fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, DownloadResponse{
					Errors: []string{"INTERNAL"},
				})
			}

			return
		}

		// body.AuthorizationCode != "" so this should be a successful login attempt, but just in case
		if len(attemptRow.Code) != core.AUTH_CODE_BYTE_LENGTH {
			fmt.Printf("warning: attemptRow.Code was the wrong length! this shouldn't happen. len(attemptRow.Code): %v\n", len(attemptRow.Code))
			ctx.JSON(http.StatusInternalServerError, DownloadResponse{
				Errors: []string{"INTERNAL"},
			})
			return
		}
		if clock.Now().UTC().Before(attemptRow.CodeValidFrom) {
			ctx.JSON(http.StatusConflict, DownloadResponse{
				Errors:                   []string{"CODE_NOT_VALID_YET"},
				AuthorizationCodeValidAt: &attemptRow.CodeValidFrom,
			})
			return
		}

		userRow, err := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Select(
				user.FieldUsername,
				// Contacts aren't needed

				user.FieldContent,
				user.FieldFileName,
				user.FieldMime,
				user.FieldNonce,
				user.FieldKeySalt,
				user.FieldPasswordHash,
				user.FieldPasswordSalt,
				user.FieldHashTime,
				user.FieldHashMemory,
				user.FieldHashKeyLen,
			).
			Only(context.Background())
		if err != nil {
			fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, DownloadResponse{
				Errors: []string{"INTERNAL"},
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
				sendUnauthorizedError(ctx)
			} else {
				fmt.Printf("warning: an error occurred while decrypting user data:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, DownloadResponse{
					Errors: []string{"INTERNAL"},
				})
			}
			return
		}

		ctx.JSON(http.StatusOK, DownloadResponse{
			Errors:                   []string{},
			AuthorizationCodeValidAt: nil,
			Content:                  decrypted,
			Filename:                 userRow.FileName,
			Mime:                     userRow.Mime,
		})

		// TODO: log this event to database
	}
}
