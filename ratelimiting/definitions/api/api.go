package api

import (
	"time"

	"github.com/NicoClack/cryptic-stash/ratelimiting"
)

func Register(group *ratelimiting.Group) {
	group.Register("", -1, 500, 10*time.Minute)
}
