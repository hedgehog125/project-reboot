package users

import (
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	group.POST("/get-authorization-code", GetAuthorizationCode(group.App))
	group.POST("/download", Download(group.App))
	group.POST("/self-lock", SelfLock(group.App))
}
