package dbcommon

import (
	"context"
	"sync"
	"testing"

	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/stretchr/testify/require"
)

func TestWithReadTx_allowsConcurrentReads(t *testing.T) {
	// TODO
}
func TestWithWriteTx_nestedTransactions_returnsError(t *testing.T) {
	// TODO
}
func TestWithWriteTx_supports25ConcurrentWrites(t *testing.T) {
	// SQLite isn't suitable if the program has many more concurrent writes than this
	t.Parallel()
	JOB_COUNT := 25
	db := testcommon.CreateDB()
	defer db.Shutdown()

	var wg sync.WaitGroup
	createJob := func() {
		wg.Add(1)
		defer wg.Done()
		stdErr := WithWriteTx(t.Context(), db, func(tx *ent.Tx, ctx context.Context) error {
			_, stdErr := tx.Job.Create().
				SetType("test_job").
				SetVersion(1).
				SetPriority(1).
				SetWeight(1).
				SetData("").
				Save(ctx)
			return stdErr
		})
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
