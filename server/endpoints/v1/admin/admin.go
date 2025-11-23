package admin

import (
	"github.com/hedgehog125/project-reboot/server/endpoints/v1/admin/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	users.ConfigureEndpoints(group.Group("/users"))
}
