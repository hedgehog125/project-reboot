package services

import (
	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/ratelimiting"
	"github.com/NicoClack/cryptic-stash/ratelimiting/definitions"
)

type RateLimiter struct {
	App     *common.App
	Limiter *ratelimiting.Limiter
}

func NewRateLimiter(app *common.App) *RateLimiter {
	limiter := ratelimiting.NewLimiter(app)
	definitions.Register(limiter.Group(""))
	return &RateLimiter{
		App:     app,
		Limiter: limiter,
	}
}

func (service *RateLimiter) RequestSession(
	eventName string, amount int, userID string,
) (common.LimiterSession, common.WrappedError) {
	session, wrappedErr := service.Limiter.RequestSession(eventName, amount, userID)
	if session == nil { // Avoid wrapping nil sessions in a non-nil interface
		return nil, wrappedErr
	}
	return session, wrappedErr
}
func (service *RateLimiter) DeleteInactiveUsers() {
	service.Limiter.DeleteInactiveUsers()
}
