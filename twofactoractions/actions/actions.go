package actions

import (
	"github.com/hedgehog125/project-reboot/twofactoractions"
	"github.com/hedgehog125/project-reboot/twofactoractions/actions/users"
)

func RegisterActions(group *twofactoractions.RegistryGroup) {
	users.RegisterActions(group.Group("users"))
}
