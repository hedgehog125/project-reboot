package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/endpoints/users"
)

func RootRedirect() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	}
}

func ConfigureEndpoints(group *gin.RouterGroup) {
	group.GET("/", RootRedirect())
	users.ConfigureEndpoints(group.Group("/api/v1/users"))
}
