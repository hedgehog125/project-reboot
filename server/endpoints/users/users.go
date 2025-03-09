package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *gin.RouterGroup, app *servercommon.ServerApp) {
	group.POST("/download", Download(app))
	group.POST("/get-authorization-code", GetAuthorizationCode(app))
	group.POST("/register-or-update", app.AdminMiddleware, RegisterOrUpdate(app))
	group.POST("/set-user-contacts", app.AdminMiddleware, SetContacts(app))
}
