package services_test

import (
	"sync"
	"testing"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/common/testcommon"
	"github.com/NicoClack/cryptic-stash/common/testcommon/mocks"
	"github.com/NicoClack/cryptic-stash/services"
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

func TestSchedulerShutdown_NoOpWhenNotStarted(t *testing.T) {
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

	select {
	case <-common.NewCallbackChannel(app.Scheduler.Shutdown):
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Scheduler Shutdown blocked when service was not started; expected no-op")
	}
}
