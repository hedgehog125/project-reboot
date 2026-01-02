package ratelimiting

import (
	"time"

	"github.com/NicoClack/cryptic-stash/backend/common"
)

type Group struct {
	Limiter *Limiter
	Path    string
}

func (limiter *Limiter) Group(relativePath string) *Group {
	return &Group{
		Limiter: limiter,
		Path:    relativePath,
	}
}

func (group *Group) Group(relativePath string) *Group {
	return &Group{
		Limiter: group.Limiter,
		Path:    common.JoinPaths(group.Path, relativePath),
	}
}

func (group *Group) Register(eventName string, globalMax, userMax int, resetDuration time.Duration) {
	fullName := common.JoinPaths(group.Path, eventName)
	group.Limiter.Register(fullName, globalMax, userMax, resetDuration)
}
