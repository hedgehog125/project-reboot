package v1

import (
	"github.com/NicoClack/cryptic-stash/server/endpoints/v1/admin"
	"github.com/NicoClack/cryptic-stash/server/endpoints/v1/twofactoractions"
	"github.com/NicoClack/cryptic-stash/server/endpoints/v1/users"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	users.ConfigureEndpoints(group.Group("/users"))
	twofactoractions.ConfigureEndpoints(group.Group("/two-factor-actions"))

	adminGroup := group.Group("/admin")
	adminGroup.Use(group.App.AdminMiddleware)
	admin.ConfigureEndpoints(adminGroup)
}
