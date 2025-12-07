package services_test

import (
	"sync"
	"testing"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/services"
)

func TestLoggerShutdown_HandlesConcurrentCalls(t *testing.T) {
	t.Parallel()

	app := &common.App{
		Env:      testcommon.DefaultEnv(),
		Database: testcommon.CreateDB(),
	}
	app.Database.Start()
	t.Cleanup(app.Database.Shutdown)
	app.Logger = services.NewLogger(app)
	app.Logger.Start()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(app.Logger.Shutdown)
	}
	wg.Wait()
}
