package definitions

import (
	"github.com/hedgehog125/project-reboot/ratelimiting"
	"github.com/hedgehog125/project-reboot/ratelimiting/definitions/api"
)

func Register(group *ratelimiting.Group) {
	api.Register(group.Group("api"))
}
