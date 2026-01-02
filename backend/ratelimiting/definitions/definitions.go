package definitions

import (
	"github.com/NicoClack/cryptic-stash/backend/ratelimiting"
	"github.com/NicoClack/cryptic-stash/backend/ratelimiting/definitions/api"
)

func Register(group *ratelimiting.Group) {
	api.Register(group.Group("api"))
	group.Register("admin-error-message", 1, -1, group.Limiter.App.Env.MIN_ADMIN_MESSAGE_GAP)
}
