package messengers

import "github.com/NicoClack/cryptic-stash/backend/server/servercommon"

func ConfigureEndpoints(group *servercommon.Group) {
	group.GET("/", ListMessengers(group.App))
	group.POST("/enable/", EnableMessenger(group.App))
	group.POST("/disable/", DisableMessenger(group.App))
}
