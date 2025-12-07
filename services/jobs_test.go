package services_test

import (
	"sync"
	"testing"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/services"
)

func TestJobsShutdown_HandlesConcurrentCalls(t *testing.T) {
	t.Parallel()

	app := &common.App{
		Env:      testcommon.DefaultEnv(),
		Database: testcommon.CreateDB(),
		Logger:   testcommon.NewTestLogger(),
	}
	app.Database.Start()
	t.Cleanup(app.Database.Shutdown)
	jobService := services.NewJobs(app)
	jobService.Start()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(jobService.Shutdown)
	}
	wg.Wait()
}
