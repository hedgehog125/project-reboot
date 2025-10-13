package jobs

import (
	"context"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestEngine_runsJob(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	defer db.Shutdown()
	app := &common.App{
		Database: db,
		Env:      testcommon.DefaultEnv(),
		Logger:   testcommon.NewTestLogger(),
		Clock:    clockwork.NewRealClock(),
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
	stdErr := dbcommon.WithWriteTx(
		t.Context(), db,
		func(tx *ent.Tx, ctx context.Context) error {
			_, commErr := engine.Enqueue("test_job_1", &body{}, ctx)
			return commErr.StandardError()
		},
	)
	require.NoError(t, stdErr)
	select {
	case <-completeJobChan:
	case <-time.After(2 * time.Second):
		t.Fatal("job did not complete in time")
	}
}

func TestEngine_retriesJob(t *testing.T) {
	t.Parallel()
	db := testcommon.CreateDB()
	defer db.Shutdown()
	app := &common.App{
		Clock:    clockwork.NewRealClock(),
		Database: db,
		Env:      testcommon.DefaultEnv(),
		Logger:   testcommon.NewTestLogger(),
	}
	type body = struct{}
	completeJobChan := make(chan struct{})
	registry := NewRegistry(app)
	attempt := 0
	registry.Register(&Definition{
		ID:      "test_job",
		Version: 1,
		Handler: func(ctx *Context) error {
			attempt++
			if attempt < 3 {
				return common.NewErrorWithCategories("temporary error").
					ConfigureRetries(2, 10*time.Millisecond, 2)
			}
			completeJobChan <- struct{}{}
			return nil
		},
		BodyType: &body{},
		Weight:   1,
	})
	engine := NewEngine(registry)
	go engine.Listen()
	defer engine.Shutdown()

	// TODO: use t.Log
	stdErr := dbcommon.WithWriteTx(
		t.Context(), db,
		func(tx *ent.Tx, ctx context.Context) error {
			_, commErr := engine.Enqueue("test_job_1", &body{}, ctx)
			return commErr.StandardError()
		},
	)
	require.NoError(t, stdErr)
	select {
	case <-completeJobChan:
	case <-time.After(2 * time.Second):
		t.Fatal("job did not complete in time")
	}
	require.Equal(t, 3, attempt)

	// TODO: not exiting gracefully
}
