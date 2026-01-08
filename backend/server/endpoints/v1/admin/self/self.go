package self

import (
	"github.com/NicoClack/cryptic-stash/backend/server/endpoints/v1/admin/self/messengers"
	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	messengers.ConfigureEndpoints(group.Group("/messengers"))
}
