package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type AdminLockPayload struct {
	Username string `binding:"required,min=1,max=32,alphanum,lowercase" json:"username"`
}
type AdminLockResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

// Admin route
func AdminLock(app *servercommon.ServerApp) gin.HandlerFunc {
	// dbClient := app.Database.Client()
	// messenger := app.Messenger

	return func(ginCtx *gin.Context) {
		body := AdminLockPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			ginCtx.Error(ctxErr)
			return
		}

		// TODO: implement
	}
}
