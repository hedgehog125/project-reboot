package users

import "github.com/hedgehog125/project-reboot/jobs"

func Register(group *jobs.RegistryGroup) {
	group.Register(TempSelfLock1(group.Registry.App))
}
