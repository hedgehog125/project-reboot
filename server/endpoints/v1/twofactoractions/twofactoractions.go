package twofactoractions

import (
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	group.POST("/:id/confirm", Confirm(group.App))
}
