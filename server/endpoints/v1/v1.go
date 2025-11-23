package v1

import (
	"github.com/hedgehog125/project-reboot/server/endpoints/v1/admin"
	"github.com/hedgehog125/project-reboot/server/endpoints/v1/twofactoractions"
	"github.com/hedgehog125/project-reboot/server/endpoints/v1/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	users.ConfigureEndpoints(group.Group("/users"))
	twofactoractions.ConfigureEndpoints(group.Group("/two-factor-actions"))

	adminGroup := group.Group("/admin")
	adminGroup.Use(group.App.AdminMiddleware)
	admin.ConfigureEndpoints(adminGroup)
}
