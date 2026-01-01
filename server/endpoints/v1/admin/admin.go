package admin

import (
	"github.com/NicoClack/cryptic-stash/server/endpoints/v1/admin/users"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	// /login is registered in v1.go since it's unauthenticated
	users.ConfigureEndpoints(group.Group("/users"))
}
