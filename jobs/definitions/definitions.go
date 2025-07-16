package definitions

import (
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/jobs/definitions/users"
)

func Register(group *jobs.RegistryGroup) {
	users.Register(group.Group("users"))
}
