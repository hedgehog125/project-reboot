package ratelimiting_test

import (
	"sync"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ratelimiting"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestRequestSession(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	limiter := ratelimiting.NewLimiter(clock)
	limiter.Register("api", -1, 50, 15*time.Minute)
	limiter.Register("hash-password", 51, 25, 15*time.Minute)

	makeGeneralRequest := func(user string) *common.Error {
		_, commErr := limiter.RequestSession("api", 1, user)
		return commErr
	}
	makeHashRequest := func(user string, shouldHashingSucceed bool) *common.Error {
		generalSession, commErr := limiter.RequestSession("api", 1, user)
		if commErr != nil {
			return commErr
		}
		hashSession, commErr := limiter.RequestSession("hash-password", 1, user)
		if commErr != nil {
			generalSession.Cancel()
			return commErr
		}
		if !shouldHashingSucceed {
			// Probably shouldn't refund in this situation, but useful for testing
			generalSession.Cancel()
			hashSession.Cancel()
			return common.NewErrorWithCategories("hashing failed")
		}
		return nil
	}
	testRateLimits := func() {
		var wg sync.WaitGroup
		for range 25 {
			wg.Go(func() {
				require.NoError(t, makeHashRequest("user1", true).StandardError())
			})
		}
		wg.Wait()
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded)
		session, commErr := limiter.RequestSession("api", 1, "user1")
		require.NoError(t, commErr.StandardError())
		session.Cancel()

		require.NoError(t, makeHashRequest("user2", true).StandardError())

		for range 25 {
			wg.Go(func() {
				require.NoError(t, makeGeneralRequest("user1").StandardError())
			})
		}
		for range 24 {
			wg.Go(func() {
				require.NoError(t, makeHashRequest("user2", true).StandardError())
			})
		}
		wg.Wait()
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded)
		require.ErrorIs(t, makeHashRequest("user2", true), ratelimiting.ErrUserRateLimitExceeded)
		require.NoError(t, makeHashRequest("user3", true).StandardError())                        // Reach the global limit for hash-password
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded) // There's no global limit for api
		require.ErrorIs(t, makeHashRequest("user2", true), ratelimiting.ErrGlobalRateLimitExceeded)
	}

	testRateLimits()
	clock.Advance(15 * time.Minute)
	testRateLimits() // Should behave the same
}
