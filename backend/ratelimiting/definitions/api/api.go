package api

import (
	"time"

	"github.com/NicoClack/cryptic-stash/backend/ratelimiting"
)

func Register(group *ratelimiting.Group) {
	group.Register("", -1, 500, 10*time.Minute)
}
