package users

import "github.com/gin-gonic/gin"

func ConfigureEndpoints(group *gin.RouterGroup) {
	group.POST("/download", Download())
	group.POST("/get-authorization-code", GetAuthorizationCode())
	group.POST("/register-or-update", adminMiddleware, RegisterOrUpdate())
	group.POST("/set-user-contacts", adminMiddleware, SetContacts())
}
