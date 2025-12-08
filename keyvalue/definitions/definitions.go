package definitions

import (
	"time"

	"github.com/NicoClack/cryptic-stash/keyvalue"
)

func Register(group *keyvalue.RegistryGroup) {
	group.Register(LastCrashSignal())
}

func LastCrashSignal() *keyvalue.Definition {
	return &keyvalue.Definition{
		Name: "LAST_CRASH_SIGNAL",
		Type: time.Time{},
		Init: func() any {
			return time.Time{} // Allow a crash signal as soon as the server first starts
		},
	}
}
