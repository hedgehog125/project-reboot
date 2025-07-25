package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
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

	return servercommon.WithTx(app, func(ctx *gin.Context, tx *ent.Tx) error {
		body := SelfLockPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			return ctxErr
		}
		until := clock.Now().Add(
			min(
				body.Until.Sub(clock.Now()), // Convert to duration
				MAX_SELF_LOCK_DURATION,
			),
		)

		userRow, stdErr := tx.User.Query().
			Where(user.Username(body.Username)).
			Only(ctx)
		if stdErr != nil {
			return servercommon.SendUnauthorizedIfNotFound(stdErr)
		}
		// TODO: check the user isn't locked

		encryptionKey := core.HashPassword(
			body.Password,
			userRow.KeySalt,
			&common.PasswordHashSettings{
				Time:    userRow.HashTime,
				Memory:  userRow.HashMemory,
				Threads: userRow.HashThreads,
			},
		)
		_, commErr := core.Decrypt(userRow.Content, encryptionKey, userRow.Nonce)
		if commErr != nil {
			return servercommon.NewUnauthorizedError()
		}

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

		commErr = app.Messengers.SendUsingAll(common.Message{
			Type: common.Message2FA,
			Code: code,
			User: (
			//exhaustruct:enforce
			&common.UserContacts{
				Username:       body.Username,
				AlertDiscordId: userRow.AlertDiscordId,
				AlertEmail:     userRow.AlertEmail,
			}),
		})
		if commErr != nil {
			return commErr
		}

		// TODO: log these errors

		ctx.JSON(http.StatusCreated, SelfLockResponse{
			Errors:            []servercommon.ErrorDetail{},
			TwoFactorActionID: actionID.String(),
		})
		return nil
	})
}
