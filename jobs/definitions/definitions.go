package definitions

import (
	"github.com/NicoClack/cryptic-stash/jobs"
	"github.com/NicoClack/cryptic-stash/jobs/definitions/users"
)

func Register(group *jobs.RegistryGroup) {
	users.Register(group.Group("users"))
}
