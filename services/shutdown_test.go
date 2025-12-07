package services_test

import (
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
)

func TestShutdown_HandlesConcurrentCalls(t *testing.T) {
	t.Parallel()

	app := &common.App{
		Clock:  clockwork.NewRealClock(),
		Logger: testcommon.NewTestLogger(),
	}
	shutdownService := services.NewShutdown(
		app,
		services.NewShutdownTask(func() {
			app.Logger.Info("doing some shutdown task...")
			app.Clock.Sleep(10 * time.Millisecond)
		}, false),
	)
	app.ShutdownService = shutdownService

	for range 100 {
		go app.Shutdown("")
	}
	shutdownService.ListenForShutdownCall()
}
