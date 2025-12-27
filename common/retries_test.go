package common_test

import (
	"testing"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/stretchr/testify/require"
)

func TestCalculateBackoff(t *testing.T) {
	t.Parallel()
	for range 1000 { // I love testing randomness :grimacing:
		shortBackoff1 := common.CalculateBackoff(0, 100*time.Millisecond, 2)
		shortBackoff2 := common.CalculateBackoff(1, 100*time.Millisecond, 2)
		longBackoff1 := common.CalculateBackoff(0, 30*time.Second, 2)
		longBackoff2 := common.CalculateBackoff(1, 30*time.Second, 2)

		require.Greater(t, shortBackoff1, 94*time.Millisecond)
		require.Less(t, shortBackoff1, 106*time.Millisecond)
		require.Greater(t, shortBackoff2, 189*time.Millisecond)
		require.Less(t, shortBackoff2, 211*time.Millisecond)

		maxDiff := 501 * time.Millisecond
		require.Greater(t, longBackoff1, (30*time.Second)-maxDiff)
		require.Less(t, longBackoff1, (30*time.Second)+maxDiff)
		require.Greater(t, longBackoff2, (60*time.Second)-maxDiff)
		require.Less(t, longBackoff2, (60*time.Second)+maxDiff)
	}
}
