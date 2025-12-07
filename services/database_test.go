package services_test

import (
	"sync"
	"testing"
	"time"

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

func TestDatabaseShutdown_NoOpWhenNotStarted(t *testing.T) {
	t.Parallel()

	env := testcommon.DefaultEnv()
	env.MOUNT_PATH = t.TempDir()
	db := services.NewDatabase(&common.App{
		Env: env,
	})

	select {
	case <-common.NewCallbackChannel(db.Shutdown):
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Shutdown blocked when service was not started; expected no-op")
	}
}
