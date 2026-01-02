package users

import "github.com/NicoClack/cryptic-stash/backend/jobs"

func Register(group *jobs.RegistryGroup) {
	group.Register(TempSelfLock1(group.Registry.App))
	group.Register(TempSelfUnlock1(group.Registry.App))
}
