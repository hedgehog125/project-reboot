package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/server/endpoints/twofactoractions"
	"github.com/hedgehog125/project-reboot/server/endpoints/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func RootRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	}
}

func ConfigureEndpoints(group *gin.RouterGroup, app *servercommon.ServerApp) {
	group.GET("/", RootRedirect())
	users.ConfigureEndpoints(group.Group("/api/v1/users"), app)
	twofactoractions.ConfigureEndpoints(group.Group("/api/v1/two-factor-actions"), app)
}
