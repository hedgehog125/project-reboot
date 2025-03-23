package users

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type DownloadPayload struct {
	Username          string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password          string `binding:"required,min=8,max=256"                   json:"password"`
	AuthorizationCode string `binding:"required,min=128,max=256"                 json:"authorizationCode"`
	// ^ I think the length can vary because of the base64 encoding?
}

type DownloadResponse struct {
	Errors                   []string   `binding:"required"              json:"errors"`
	AuthorizationCodeValidAt *time.Time `json:"authorizationCodeValidAt"`
	Content                  []byte     `json:"content"`
	Filename                 string     `json:"filename"`
	Mime                     string     `json:"mime"`
}

func Download(app *servercommon.ServerApp) gin.HandlerFunc {
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

		sessionRow, err := dbClient.Session.Query().
			Where(session.And(session.HasUserWith(user.Username(body.Username)), session.Code(givenAuthCodeBytes))).
			Select(session.FieldCode, session.FieldCodeValidFrom).
			First(context.Background())
		if err != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(err))
			return
		}

		if clock.Now().UTC().Before(sessionRow.CodeValidFrom) {
			ctx.JSON(http.StatusConflict, DownloadResponse{
				Errors:                   []string{"CODE_NOT_VALID_YET"},
				AuthorizationCodeValidAt: &sessionRow.CodeValidFrom,
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
			ctx.Error(err)
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
			ctx.Error(servercommon.ExpectError(err, core.ErrIncorrectPassword, http.StatusUnauthorized, "UNAUTHORIZED"))
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
