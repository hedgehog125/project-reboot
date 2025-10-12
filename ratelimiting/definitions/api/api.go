package api

import (
	"time"

	"github.com/hedgehog125/project-reboot/ratelimiting"
)

func Register(group *ratelimiting.Group) {
	group.Register("", -1, 500, 10*time.Minute)
}
