package dbcommon_test

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/common/testcommon"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/ent/job"
	"github.com/stretchr/testify/require"
)

func TestWithReadTx_allowsConcurrentReads(t *testing.T) {
	// TODO
}
func TestWithWriteTx_nestedTransactions_returnsError(t *testing.T) {
	// TODO
}
func TestWithWriteTx_supports50ConcurrentWrites(t *testing.T) {
	// SQLite isn't suitable if the program has many more concurrent writes than this
	t.Parallel()
	JOB_COUNT := 50
	db := testcommon.CreateDB() // TODO: use a disk database to more accurately measure performance
	defer db.Shutdown()

	var wg sync.WaitGroup
	createJob := func() {
		wg.Add(1)
		defer wg.Done()
		stdErr := dbcommon.WithWriteTx(
			t.Context(), db,
			func(tx *ent.Tx, ctx context.Context) error {
				_, stdErr := tx.Job.Create().
					SetType("test_job").
					SetCreatedAt(time.Now()).
					SetDueAt(time.Now()).
					SetOriginallyDueAt(time.Now()).
					SetVersion(1).
					SetPriority(1).
					SetWeight(1).
					SetBody(json.RawMessage("{}")).
					Save(ctx)
				if stdErr != nil {
					return common.ErrWrapperDatabase.Wrap(stdErr)
				}
				return nil
			},
		)
		require.NoError(t, stdErr)
	}
	for range JOB_COUNT {
		go createJob()
	}
	wg.Wait()
	count, stdErr := db.Client().Job.Query().Count(t.Context())
	require.NoError(t, stdErr)
	require.Equal(t, JOB_COUNT, count)
}
func TestWithWriteTx_supports25CollidingIncrements(t *testing.T) {
	t.Parallel()
	INCREMENT_COUNT := 25
	db := testcommon.CreateDB()
	defer db.Shutdown()

	stdErr := db.Client().Job.Create().
		SetType("counter").
		SetCreatedAt(time.Now()).
		SetDueAt(time.Now()).
		SetOriginallyDueAt(time.Now()).
		SetVersion(1).
		SetPriority(1).
		SetWeight(1).
		SetBody(json.RawMessage(`{"count":0}`)).
		Exec(t.Context())
	require.NoError(t, stdErr)

	var errCount atomic.Int32
	var wg sync.WaitGroup
	for range INCREMENT_COUNT {
		wg.Go(func() {
			stdErr := dbcommon.WithWriteTx(
				t.Context(), db,
				func(tx *ent.Tx, ctx context.Context) error {
					job, stdErr := tx.Job.Query().Where(job.TypeEQ("counter")).Only(ctx)
					if stdErr != nil {
						errCount.Add(1)
						return common.ErrWrapperDatabase.Wrap(stdErr)
					}
					// Extend the window when this transaction hasn't got a write lock
					time.Sleep(10 * time.Millisecond)

					var body struct {
						Count int `json:"count"`
					}
					stdErr = json.Unmarshal(job.Body, &body)
					if stdErr != nil {
						errCount.Add(1)
						return stdErr
					}
					body.Count++
					newBody, stdErr := json.Marshal(body)
					if stdErr != nil {
						errCount.Add(1)
						return stdErr
					}
					stdErr = job.Update().SetBody(json.RawMessage(newBody)).Exec(ctx)
					if stdErr != nil {
						errCount.Add(1)
						return common.ErrWrapperDatabase.Wrap(stdErr)
					}
					return nil
				},
			)
			require.NoError(t, stdErr)
		})
	}
	wg.Wait()

	jobOb, err := db.Client().Job.Query().Where(job.TypeEQ("counter")).Only(t.Context())
	require.NoError(t, err)
	var body struct {
		Count int `json:"count"`
	}
	stdErr = json.Unmarshal(jobOb.Body, &body)
	require.NoError(t, stdErr)
	require.Equal(t, INCREMENT_COUNT, body.Count)
	// Expect at least a few errors to have been retried (it should be way more than this on average)
	require.GreaterOrEqual(t, errCount.Load(), int32(5))
}
