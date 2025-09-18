package schedulers

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/periodictask"
)

type DelayFuncContext struct {
	App     *common.App
	LastRan time.Time
	Context context.Context
}

// You should return the initial value for lastRan if lastRan is zero
type DelayFunc = func(delayCtx *DelayFuncContext) (nextRun time.Time, commit CommitDelayFunc)
type CommitDelayFunc = func(runTime time.Time, ctx context.Context)

func SimpleFixedInterval(interval time.Duration) DelayFunc {
	return func(delayCtx *DelayFuncContext) (time.Time, func(runTime time.Time, ctx context.Context)) {
		lastRan := delayCtx.LastRan
		if lastRan.IsZero() {
			lastRan = delayCtx.App.Clock.Now()
		}
		return lastRan.Add(interval), func(runTime time.Time, ctx context.Context) {}
	}
}
func PersistentFixedInterval(periodicTaskName string, interval time.Duration) DelayFunc {
	periodicTaskID := 0
	return func(delayCtx *DelayFuncContext) (time.Time, CommitDelayFunc) {
		lastRan := delayCtx.LastRan
		commit := func(runTime time.Time, ctx context.Context) {
			if periodicTaskID == 0 {
				periodicTask, stdErr := dbcommon.WithReadWriteTx(
					ctx,
					delayCtx.App.Database,
					func(tx *ent.Tx, ctx context.Context) (*ent.PeriodicTask, error) {
						return tx.PeriodicTask.Create().
							SetName(periodicTaskName).
							SetLastRan(runTime).
							Save(ctx)
					},
				)
				if stdErr == nil {
					periodicTaskID = periodicTask.ID
				} else {
					delayCtx.App.Logger.Error(
						"unable to create initial PeriodicTask object",
						"error", stdErr,
						"periodicTaskName", periodicTaskName,
					)
				}
			} else {
				stdErr := dbcommon.WithWriteTx(
					ctx,
					delayCtx.App.Database,
					func(tx *ent.Tx, ctx context.Context) error {
						return tx.PeriodicTask.UpdateOneID(periodicTaskID).
							SetLastRan(runTime).
							Exec(ctx)
					},
				)
				if stdErr != nil {
					delayCtx.App.Logger.Error(
						"unable to update lastRun of PeriodicTask object",
						"error", stdErr,
						"periodicTaskName", periodicTaskName,
						"periodicTaskID", periodicTaskID,
					)
				}
			}
		}

		if lastRan.IsZero() {
			periodicTask, stdErr := dbcommon.WithReadTx(
				delayCtx.Context,
				delayCtx.App.Database,
				func(tx *ent.Tx, ctx context.Context) (*ent.PeriodicTask, error) {
					return tx.PeriodicTask.Query().
						Where(periodictask.Name(periodicTaskName)).
						Only(ctx)
				},
			)
			if stdErr == nil {
				periodicTaskID = periodicTask.ID
				lastRan = periodicTask.LastRan
			} else {
				if !ent.IsNotFound(stdErr) {
					delayCtx.App.Logger.Error(
						"couldn't read PeriodicTask object on startup, assuming this is the first run",
						"error", stdErr,
						"periodicTaskName", periodicTaskName,
					)
				}
				return delayCtx.App.Clock.Now(), commit
			}
		}
		return lastRan.Add(interval), commit
	}
}

// TODO: add function to run at regular times, with an offset argument (e.g run every day at 9am)
