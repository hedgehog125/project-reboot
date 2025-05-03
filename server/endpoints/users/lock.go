package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

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
		body := SetContactsPayload{}
		if err := ctx.BindJSON(&body); err != nil {
			return
		}

		// TODO: implement
	}
}

type LockTemporarilyPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
}
type LockTemporarilyResponse struct {
	Errors            []string `binding:"required"  json:"errors"`
	TwoFactorActionID string   `json:"twoFactorActionID"`
}

func LockTemporarily(app *servercommon.ServerApp) gin.HandlerFunc {
	dbClient := app.App.Database.Client()
	clock := app.App.Clock

	return func(ctx *gin.Context) {
		body := SetContactsPayload{}
		if err := ctx.BindJSON(&body); err != nil {
			return
		}

		// TODO: auth!

		actionID, _, err := twofactoractions.Create(
			"TEMP_SELF_LOCK", 1,
			clock.Now().Add(twofactoractions.DEFAULT_CODE_LIFETIME),
			twofactoractions.TempSelfLock1{},
			dbClient,
		)
		if err != nil {
			// TODO: categorise the database errors properly
			ctx.Error(err)
			return
		}

		// TODO: send code

		ctx.JSON(http.StatusCreated, LockTemporarilyResponse{
			Errors:            []string{},
			TwoFactorActionID: actionID.String(),
		})
	}
}
