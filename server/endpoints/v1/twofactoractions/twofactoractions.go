package twofactoractions

import (
	"github.com/NicoClack/cryptic-stash/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	group.POST("/:id/confirm", Confirm(group.App))
}
