package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ratelimiting"
	"github.com/hedgehog125/project-reboot/ratelimiting/definitions"
)

type RateLimiter struct {
	App     *common.App
	Limiter *ratelimiting.Limiter
}

func NewRateLimiter(app *common.App) *RateLimiter {
	limiter := ratelimiting.NewLimiter(app.Clock)
	definitions.Register(limiter.Group(""))
	return &RateLimiter{
		App:     app,
		Limiter: limiter,
	}
}

func (service *RateLimiter) RequestSession(
	eventName string, amount int, userID string,
) (common.LimiterSession, *common.Error) {
	session, commErr := service.Limiter.RequestSession(eventName, amount, userID)
	if session == nil { // Avoid wrapping nil sessions in a non-nil interface
		return nil, commErr
	}
	return session, commErr
}
