package users

import "github.com/NicoClack/cryptic-stash/jobs"

func Register(group *jobs.RegistryGroup) {
	group.Register(TempSelfLock1(group.Registry.App))
	group.Register(TempSelfUnlock1(group.Registry.App))
}
