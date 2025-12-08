package ratelimiting_test

import (
	"sync"
	"testing"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/ratelimiting"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestRequestSession(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	limiter := ratelimiting.NewLimiter(&common.App{
		Clock: clock,
	})
	limiter.Register("api", -1, 50, 15*time.Minute)
	limiter.Register("hash-password", 51, 25, 15*time.Minute)

	makeGeneralRequest := func(user string) common.WrappedError {
		_, wrappedErr := limiter.RequestSession("api", 1, user)
		return wrappedErr
	}
	makeHashRequest := func(user string, shouldHashingSucceed bool) common.WrappedError {
		generalSession, wrappedErr := limiter.RequestSession("api", 1, user)
		if wrappedErr != nil {
			return wrappedErr
		}
		hashSession, wrappedErr := limiter.RequestSession("hash-password", 1, user)
		if wrappedErr != nil {
			generalSession.Cancel()
			return wrappedErr
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
				require.NoError(t, makeHashRequest("user1", true))
			})
		}
		wg.Wait()
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded)
		session, wrappedErr := limiter.RequestSession("api", 1, "user1")
		require.NoError(t, wrappedErr)
		session.Cancel()

		require.NoError(t, makeHashRequest("user2", true))

		for range 25 {
			wg.Go(func() {
				require.NoError(t, makeGeneralRequest("user1"))
			})
		}
		for range 24 {
			wg.Go(func() {
				require.NoError(t, makeHashRequest("user2", true))
			})
		}
		wg.Wait()
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded)
		require.ErrorIs(t, makeHashRequest("user2", true), ratelimiting.ErrUserRateLimitExceeded)
		require.NoError(t, makeHashRequest("user3", true))                                        // Reach the global limit for hash-password
		require.ErrorIs(t, makeHashRequest("user1", true), ratelimiting.ErrUserRateLimitExceeded) // There's no global limit for api
		require.ErrorIs(t, makeHashRequest("user2", true), ratelimiting.ErrGlobalRateLimitExceeded)
	}

	testRateLimits()
	clock.Advance(15 * time.Minute)
	testRateLimits() // Should behave the same
}
