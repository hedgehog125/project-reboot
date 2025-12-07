package services_test

import (
	"sync"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/common/testcommon/mocks"
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
)

func TestSchedulerShutdown_HandlesConcurrentCalls(t *testing.T) {
	t.Parallel()

	app := &common.App{
		Clock:            clockwork.NewRealClock(),
		Env:              testcommon.DefaultEnv(),
		Database:         testcommon.CreateDB(),
		Logger:           testcommon.NewTestLogger(),
		Core:             mocks.NewEmptyCoreService(),
		TwoFactorActions: mocks.NewEmptyTwoFactorActionService(),
		RateLimiter:      mocks.NewEmptyRateLimiterService(),
	}
	app.Database.Start()
	t.Cleanup(app.Database.Shutdown)
	app.Scheduler = services.NewScheduler(app)
	app.Scheduler.Start()

	time.Sleep(10 * time.Millisecond) // Ensure the tasks have had a chance to start

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(app.Scheduler.Shutdown)
	}
	wg.Wait()
}
