package users

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
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
	Errors                   []servercommon.ErrorDetail `binding:"required" json:"errors"`
	AuthorizationCodeValidAt *time.Time                 `json:"authorizationCodeValidAt"`
	Content                  []byte                     `json:"content"`
	Filename                 string                     `json:"filename"`
	Mime                     string                     `json:"mime"`
}

func Download(app *servercommon.ServerApp) gin.HandlerFunc {
	clock := app.Clock

	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := DownloadPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}
		givenAuthCodeBytes, ctxErr := servercommon.DecodeBase64(body.AuthorizationCode)
		if ctxErr != nil {
			return ctxErr
		}

		userRow, stdErr := dbcommon.WithReadTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
			sessionRow, stdErr := tx.Session.Query().
				Where(session.And(session.HasUserWith(user.Username(body.Username)), session.Code(givenAuthCodeBytes))).
				Select(session.FieldCode, session.FieldCodeValidFrom).
				First(ctx)
			if stdErr != nil {
				return nil, servercommon.SendUnauthorizedIfNotFound(
					common.ErrWrapperDatabase.Wrap(stdErr),
				)
			}

			if clock.Now().Before(sessionRow.CodeValidFrom) {
				ginCtx.JSON(http.StatusConflict, DownloadResponse{
					Errors: []servercommon.ErrorDetail{
						{
							Message: "authorization code is not valid yet",
							Code:    "CODE_NOT_VALID_YET",
						},
					},
					AuthorizationCodeValidAt: &sessionRow.CodeValidFrom,
				})
				return nil, servercommon.ErrCancelTransaction.Clone()
			}

			userRow, stdErr := tx.User.Query().
				Where(user.Username(body.Username)).
				Select(
					user.FieldUsername,
					// Contacts aren't needed

					user.FieldContent,
					user.FieldFileName,
					user.FieldMime,
					user.FieldNonce,
					user.FieldKeySalt,
					user.FieldHashTime,
					user.FieldHashMemory,
					user.FieldHashThreads,
				).
				Only(ctx)
			if stdErr != nil {
				return nil, common.ErrWrapperDatabase.Wrap(stdErr)
			}
			return userRow, nil
		})
		if stdErr != nil {
			return stdErr
		}

		encryptionKey := core.HashPassword(
			body.Password,
			userRow.KeySalt,
			&common.PasswordHashSettings{
				Time:    userRow.HashTime,
				Memory:  userRow.HashMemory,
				Threads: userRow.HashThreads,
			},
		)
		decrypted, commErr := core.Decrypt(userRow.Content, encryptionKey, userRow.Nonce)
		if commErr != nil {
			return servercommon.ExpectError(
				commErr, core.ErrIncorrectPassword,
				http.StatusUnauthorized, nil,
			)
		}

		ginCtx.JSON(http.StatusOK, DownloadResponse{
			Errors:                   []servercommon.ErrorDetail{},
			AuthorizationCodeValidAt: nil,
			Content:                  decrypted,
			Filename:                 userRow.FileName,
			Mime:                     userRow.Mime,
		})
		return nil

		// TODO: log this event to database
		// TODO: reduce session expiry to 1 hour
		// TODO: notify user in the background
	})
}
