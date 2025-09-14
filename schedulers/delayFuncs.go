package schedulers

import (
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

// Note: the time.Time value in the channel is not used
type DelayFunc func(lastRan time.Time, app *common.App) time.Duration

func FixedInterval(interval time.Duration) DelayFunc {
	return func(lastRan time.Time, app *common.App) time.Duration {
		return interval - app.Clock.Since(lastRan)
	}
}

// TODO: add function to run at regular times, with an offset argument (e.g run every day at 9am)
