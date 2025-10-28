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
	Errors                      []servercommon.ErrorDetail `binding:"required" json:"errors"`
	AuthorizationCodeValidFrom  *time.Time                 `json:"authorizationCodeValidFrom"`
	AuthorizationCodeValidUntil *time.Time                 `json:"authorizationCodeValidUntil"`
	Content                     []byte                     `json:"content"`
	Filename                    string                     `json:"filename"`
	Mime                        string                     `json:"mime"`
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

		sessionOb, stdErr := dbcommon.WithReadTx(
			ginCtx, app.Database,
			func(tx *ent.Tx, ctx context.Context) (*ent.Session, error) {
				sessionOb, stdErr := tx.Session.Query().
					Where(session.And(session.HasUserWith(user.Username(body.Username)), session.Code(givenAuthCodeBytes))).
					WithUser().
					First(ctx)
				if stdErr != nil {
					return nil, servercommon.SendUnauthorizedIfNotFound(
						common.ErrWrapperDatabase.Wrap(stdErr),
					)
				}
				if clock.Now().After(sessionOb.ValidUntil) {
					stdErr := tx.Session.DeleteOneID(sessionOb.ID).Exec(ctx)
					if stdErr != nil {
						servercommon.GetLogger(ginCtx).Error(
							"unable to delete expired session",
							"error",
							stdErr,
						)
					}
					return nil, servercommon.NewUnauthorizedError()
				}
				return sessionOb, nil
			},
		)
		if stdErr != nil {
			return stdErr
		}
		if clock.Now().Before(sessionOb.ValidFrom) {
			ginCtx.JSON(http.StatusBadRequest, DownloadResponse{
				Errors: []servercommon.ErrorDetail{
					{
						Message: "authorization code is not valid yet",
						Code:    "CODE_NOT_VALID_YET",
					},
				},
				AuthorizationCodeValidFrom:  &sessionOb.ValidFrom,
				AuthorizationCodeValidUntil: &sessionOb.ValidUntil,
			})
			return nil
		}

		userOb := sessionOb.Edges.User
		encryptionKey := app.Core.HashPassword(
			body.Password,
			userOb.KeySalt,
			&common.PasswordHashSettings{
				Time:    userOb.HashTime,
				Memory:  userOb.HashMemory,
				Threads: userOb.HashThreads,
			},
		)
		decrypted, commErr := app.Core.Decrypt(userOb.Content, encryptionKey, userOb.Nonce)
		if commErr != nil {
			return servercommon.ExpectError(
				commErr, core.ErrIncorrectPassword,
				http.StatusUnauthorized, nil,
			)
		}

		stdErr = dbcommon.WithWriteTx(
			ginCtx, app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				stdErr := tx.Session.UpdateOneID(sessionOb.ID).
					SetValidUntil(clock.Now().Add(app.Env.USED_AUTH_CODE_VALID_FOR)).
					Exec(ctx)
				if stdErr != nil {
					return stdErr
				}
				_, _, commErr := app.Messengers.SendUsingAll(
					&common.Message{
						Type: common.MessageDownload,
						User: userOb,
					},
					ctx,
				)
				return commErr.StandardError()
			},
		)
		if stdErr != nil {
			return stdErr
		}

		ginCtx.JSON(http.StatusOK, DownloadResponse{
			Errors:                      []servercommon.ErrorDetail{},
			AuthorizationCodeValidFrom:  nil,
			AuthorizationCodeValidUntil: nil,
			Content:                     decrypted,
			Filename:                    userOb.FileName,
			Mime:                        userOb.Mime,
		})
		return nil
	})
}
