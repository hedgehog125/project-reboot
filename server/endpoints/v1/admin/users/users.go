package users

import "github.com/hedgehog125/project-reboot/server/servercommon"

func ConfigureEndpoints(group *servercommon.Group) {
	group.POST("/register-or-update", RegisterOrUpdate(group.App))
	group.POST("/set-user-contacts", SetContacts(group.App))
	group.POST("/lock", AdminLock(group.App))
	group.POST("/unlock", AdminUnlock(group.App))
}
