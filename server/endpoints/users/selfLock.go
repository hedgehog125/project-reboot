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
	"github.com/hedgehog125/project-reboot/ent/user"
	userjobs "github.com/hedgehog125/project-reboot/jobs/definitions/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

const MAX_SELF_LOCK_DURATION = 14 * (24 * time.Hour)

type SelfLockPayload struct {
	Username string               `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string               `binding:"required,min=8,max=256"                   json:"password"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}
type SelfLockResponse struct {
	Errors            []servercommon.ErrorDetail `binding:"required" json:"errors"`
	TwoFactorActionID string                     `json:"twoFactorActionID"`
}

func SelfLock(app *servercommon.ServerApp) gin.HandlerFunc {
	clock := app.Clock

	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := SelfLockPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}
		until := clock.Now().Add(
			min(
				body.Until.Sub(clock.Now()), // Convert to duration
				MAX_SELF_LOCK_DURATION,
			),
		)

		userOb, stdErr := dbcommon.WithReadTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
			userOb, stdErr := tx.User.Query().
				Where(user.Username(body.Username)).
				Only(ctx)
			if stdErr != nil {
				return nil, servercommon.SendUnauthorizedIfNotFound(
					common.ErrWrapperDatabase.Wrap(stdErr),
				)
			}

			return userOb, nil
		})
		if stdErr != nil {
			return stdErr
		}
		// TODO: check the user isn't locked

		encryptionKey := core.HashPassword(
			body.Password,
			userOb.KeySalt,
			&common.PasswordHashSettings{
				Time:    userOb.HashTime,
				Memory:  userOb.HashMemory,
				Threads: userOb.HashThreads,
			},
		)
		_, commErr := core.Decrypt(userOb.Content, encryptionKey, userOb.Nonce)
		if commErr != nil {
			return servercommon.NewUnauthorizedError()
		}

		return dbcommon.WithWriteTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) error {
			actionID, code, commErr := app.TwoFactorActions.Create(
				"users/TEMP_SELF_LOCK_1",
				clock.Now().Add(twofactoractions.DEFAULT_CODE_LIFETIME),
				//exhaustruct:enforce
				&userjobs.TempSelfLock1Body{
					// TODO: can this be accessed through the registry instead?
					Username: body.Username,
					Until:    common.ISOTimeString{Time: until},
				},
				ctx,
			)
			if commErr != nil {
				return commErr
			}

			_, _, commErr = app.Messengers.SendUsingAll(
				&common.Message{
					Type: common.Message2FA,
					User: userOb,
					Code: code,
				},
				ctx,
			)
			if commErr != nil {
				return commErr
			}

			// TODO: wait for job to run and return error if it fails?
			ginCtx.JSON(http.StatusCreated, SelfLockResponse{
				Errors:            []servercommon.ErrorDetail{},
				TwoFactorActionID: actionID.String(),
			})
			return nil
		})
	})
}
