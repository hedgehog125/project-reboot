package services_test

import (
	"sync"
	"testing"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/services"
)

func TestDatabaseShutdown_HandlesConcurrentCalls(t *testing.T) {
	t.Parallel()

	env := testcommon.DefaultEnv()
	env.MOUNT_PATH = t.TempDir()
	db := services.NewDatabase(&common.App{
		Env: env,
	})
	db.Start()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(db.Shutdown)
	}
	wg.Wait()
}
