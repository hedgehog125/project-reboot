package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
	useractions "github.com/hedgehog125/project-reboot/twofactoractions/actions/users"
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
	dbClient := app.Database.Client()
	clock := app.Clock
	messenger := app.Messenger

	return func(ctx *gin.Context) {
		body := SelfLockPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			ctx.Error(ctxErr)
			return
		}
		until := clock.Now().Add(
			min(
				body.Until.Sub(clock.Now()), // Convert to duration
				MAX_SELF_LOCK_DURATION,
			),
		)

		userRow, stdErr := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Only(ctx)
		if stdErr != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(stdErr))
			return
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
			ctx.Error(servercommon.NewUnauthorizedError())
			return
		}

		actionID, code, commErr := app.TwoFactorAction.Create(
			"users/TEMP_SELF_LOCK", 1,
			clock.Now().Add(twofactoractions.DEFAULT_CODE_LIFETIME),
			//exhaustruct:enforce
			&useractions.TempSelfLock1Body{
				// TODO: can this be accessed through the registry instead?
				Username: body.Username,
				Until:    common.ISOTimeString{Time: until},
			},
		)
		if commErr != nil {
			ctx.Error(commErr)
			return
		}

		errs := messenger.SendUsingAll(common.Message{
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
		if len(errs) == len(messenger.IDs()) {
			// We aren't sure if this error is the client or server's fault
			ctx.JSON(http.StatusBadRequest, SetContactsResponse{
				Errors: []servercommon.ErrorDetail{
					{
						Message: "all messages failed",
						Code:    "ALL_MESSAGES_FAILED",
					},
				},
				MessagesSent: []string{},
			})
			return
		}

		// TODO: log these errors

		ctx.JSON(http.StatusCreated, SelfLockResponse{
			Errors:            []servercommon.ErrorDetail{},
			TwoFactorActionID: actionID.String(),
		})
	}
}
