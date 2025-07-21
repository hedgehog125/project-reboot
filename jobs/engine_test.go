package jobs

import (
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/stretchr/testify/require"
)

func TestEngine_runsJob(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	defer db.Shutdown()
	app := &common.App{
		Database: db,
		Env: &common.Env{
			MAX_TOTAL_JOB_WEIGHT: 100,
		},
	}
	type body = struct{}
	completeJobChan := make(chan struct{})
	registry := NewRegistry(app)
	registry.Register(&Definition{
		ID:      "test_job",
		Version: 1,
		Handler: func(ctx *Context) error {
			completeJobChan <- struct{}{}
			return nil
		},
		BodyType: &body{},
		Weight:   1,
	})
	engine := NewEngine(registry)
	go engine.Listen()
	defer engine.Shutdown()
	_, commErr := engine.Enqueue("test_job_1", &body{}, t.Context()) // TODO: add transaction to context
	require.NoError(t, commErr.StandardError())
	select {
	case <-completeJobChan:
	case <-time.After(2 * time.Second):
		t.Fatal("job did not complete in time")
	}
}
