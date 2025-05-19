package users

import "github.com/hedgehog125/project-reboot/twofactoractions"

func RegisterActions(group *twofactoractions.RegistryGroup) {
	group.RegisterAction(TempSelfLock1(group.Registry.App))
}
