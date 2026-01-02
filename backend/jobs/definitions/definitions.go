package definitions

import (
	"github.com/NicoClack/cryptic-stash/backend/jobs"
	"github.com/NicoClack/cryptic-stash/backend/jobs/definitions/users"
)

func Register(group *jobs.RegistryGroup) {
	users.Register(group.Group("users"))
}
