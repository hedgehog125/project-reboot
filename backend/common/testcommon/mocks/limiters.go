package mocks

import "github.com/NicoClack/cryptic-stash/backend/common"

type EmptyRateLimiterService struct{}

func NewEmptyRateLimiterService() *EmptyRateLimiterService {
	return &EmptyRateLimiterService{}
}

func (m *EmptyRateLimiterService) RequestSession(
	eventName string,
	amount int,
	user string,
) (common.LimiterSession, common.WrappedError) {
	return &EmptyRateLimiterSession{}, nil
}
func (m *EmptyRateLimiterService) DeleteInactiveUsers() {
}

type EmptyRateLimiterSession struct{}

func (m *EmptyRateLimiterSession) AdjustTo(amount int) common.WrappedError {
	return nil
}
func (m *EmptyRateLimiterSession) Cancel() {
}
