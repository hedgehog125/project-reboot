package endpoints

import (
	v1 "github.com/NicoClack/cryptic-stash/server/endpoints/v1"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	v1.ConfigureEndpoints(group.Group("/api/v1"))
}
