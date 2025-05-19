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
)

const MAX_SELF_LOCK_DURATION = 14 * (24 * time.Hour)

type LockPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
}
type LockResponse struct {
	Errors []string `binding:"required"  json:"errors"`
}

// Admin route
func Lock(app *servercommon.ServerApp) gin.HandlerFunc {
	// dbClient := app.App.Database.Client()
	// messenger := app.App.Messenger

	return func(ctx *gin.Context) {
		body := LockPayload{}
		if err := ctx.BindJSON(&body); err != nil {
			return
		}

		// TODO: implement
	}
}

type LockTemporarilyPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
	Password string `binding:"required,min=8,max=256"                   json:"password"`
	Until    string `binding:"required,max=256" json:"until"`
}
type LockTemporarilyResponse struct {
	Errors            []string `binding:"required"  json:"errors"`
	TwoFactorActionID string   `json:"twoFactorActionID"`
}

func LockTemporarily(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()
	clock := app.App.Clock
	messenger := app.App.Messenger

	return func(ctx *gin.Context) {
		body := LockTemporarilyPayload{}
		if err := ctx.BindJSON(&body); err != nil {
			return
		}
		until, err := time.Parse(time.RFC3339, body.Until)
		if err != nil {
			ctx.Error(servercommon.NewBadRequestError("until", "invalid date format"))
			return
		}
		until = clock.Now().Add(
			min(
				until.Sub(clock.Now()), // Convert to duration
				MAX_SELF_LOCK_DURATION,
			),
		)

		userRow, err := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Select(
				user.FieldPasswordHash, user.FieldPasswordSalt,
				user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen,
				// Contacts
				user.FieldAlertDiscordId,
				user.FieldAlertEmail,
			).
			Only(ctx)
		if err != nil {
			ctx.Error(servercommon.SendUnauthorizedIfNotFound(err))
			return
		}

		if !core.CheckPassword(
			body.Password,
			userRow.PasswordHash,
			userRow.PasswordSalt,
			core.HashSettings{
				Time:   userRow.HashTime,
				Memory: userRow.HashMemory,
				KeyLen: userRow.HashKeyLen,
			},
		) {
			ctx.Error(servercommon.NewUnauthorizedError())
			return
		}

		actionID, code, err := twofactoractions.Create(
			"TEMP_SELF_LOCK", 1,
			clock.Now().Add(twofactoractions.DEFAULT_CODE_LIFETIME),
			//exhaustruct:enforce
			twofactoractions.TempSelfLock1{
				Username: body.Username,
				Until:    until,
			},
			dbClient,
		)
		if err != nil {
			// TODO: categorise the database errors properly
			ctx.Error(err)
			return
		}

		errs := messenger.SendUsingAll(common.Message{
			Type: common.Message2FA,
			Code: code,
			User: (
			//exhaustruct:enforce
			&common.MessageUserInfo{
				Username:       body.Username,
				AlertDiscordId: userRow.AlertDiscordId,
				AlertEmail:     userRow.AlertEmail,
			}),
		})
		messengerIDs := messenger.IDs()
		if len(errs) == len(messengerIDs) {
			// We aren't sure if this error is the client or server's fault
			ctx.JSON(http.StatusBadRequest, SetContactsResponse{
				Errors:       []string{"ALL_TEST_MESSAGES_FAILED"},
				MessagesSent: []string{},
			})
			return
		}

		// TODO: log these errors

		ctx.JSON(http.StatusCreated, LockTemporarilyResponse{
			Errors:            []string{},
			TwoFactorActionID: actionID.String(),
		})
	}
}
