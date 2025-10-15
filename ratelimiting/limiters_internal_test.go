package ratelimiting

import (
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestDeleteInactiveUsers(t *testing.T) {
	t.Parallel()
	clock := clockwork.NewFakeClock()
	limiter := NewLimiter(&common.App{
		Clock: clock,
	})
	limiter.Register("api", -1, 50, 15*time.Minute)
	_, commErr := limiter.RequestSession("api", 1, "user1")
	require.NoError(t, commErr.StandardError())
	_, ok := limiter.limits["api"].userCounters["user1"]
	require.True(t, ok)

	limiter.DeleteInactiveUsers()
	_, ok = limiter.limits["api"].userCounters["user1"]
	require.True(t, ok)

	clock.Advance(15 * time.Minute)
	limiter.DeleteInactiveUsers()
	_, ok = limiter.limits["api"].userCounters["user1"]
	require.False(t, ok)
}
