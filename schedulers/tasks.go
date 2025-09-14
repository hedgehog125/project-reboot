package schedulers

import (
	"context"
	"log"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/job"
	"github.com/hedgehog125/project-reboot/ent/periodicjob"
)

type Task struct {
	Init func(engine *Engine)
	Run  func(engine *Engine)
}

func NewTask(fn func(app *common.App), delayFunc DelayFunc) Task {
	lastRan := time.Time{}
	return Task{
		Init: func(engine *Engine) {
			lastRan = engine.App.Clock.Now()
			fn(engine.App)
		},
		Run: func(engine *Engine) {
			for {
				select {
				case <-time.After(delayFunc(lastRan, engine.App)):
				case <-engine.RequestShutdownChan:
					return
				}
				lastRan = engine.App.Clock.Now()
				fn(engine.App)
			}
		},
	}
}
func NewJobTask(versionedType string, delayFunc DelayFunc, maxConcurrentRuns int) Task {
	jobType, jobVersion, commErr := common.ParseVersionedType(versionedType)
	if commErr != nil {
		log.Fatalf("failed to parse versionedType when creating scheduler job task. error:\n%v", commErr)
	}
	if maxConcurrentRuns <= 0 && maxConcurrentRuns != -1 {
		log.Fatal("maxConcurrentRuns for scheduler job task must be >= 1 or -1 for unlimited")
	}

	lastScheduled := time.Time{}
	periodicJobID := 0
	scheduleJob := func(app *common.App, isStartup bool) time.Duration {
		sleepTime := time.Duration(0)
		// Don't update until the transaction is successful
		newLastScheduled := time.Time{}
		newPeriodicJobID := 0
		stdErr := dbcommon.WithWriteTx(
			context.TODO(), app.Database,
			func(tx *ent.Tx, ctx context.Context) error {
				newLastScheduled = lastScheduled
				newPeriodicJobID = periodicJobID

				if newPeriodicJobID != 0 {
					periodicJob, stdErr := tx.PeriodicJob.Query().
						Where(periodicjob.ID(newPeriodicJobID)).
						WithJobs(func(jobQuery *ent.JobQuery) {
							jobQuery.Select(job.FieldID, job.FieldDue)
						}).
						Only(ctx)
					if ent.IsNotFound(stdErr) {
						app.Logger.Error(
							"previously read PeriodicJob was not found when checking for running jobs, assuming none are running",
							"periodicJobID", newPeriodicJobID,
							"jobType", versionedType,
						)
						newPeriodicJobID = 0
					} else if stdErr != nil {
						return stdErr
					} else {
						if maxConcurrentRuns != -1 && len(periodicJob.Edges.Jobs) >= maxConcurrentRuns {
							if !isStartup {
								maxDue := app.Clock.Now()
								for _, job := range periodicJob.Edges.Jobs {
									if job.Due.After(maxDue) {
										maxDue = job.Due
									}
								}
								sleepTime = app.Clock.Since(maxDue) + (500 * time.Millisecond)
								app.Logger.Info(
									"maximum number of concurrent runs for periodic job reached, delaying next run",
									"periodicJobID", newPeriodicJobID,
									"jobType", versionedType,
									"maxConcurrentRuns", maxConcurrentRuns,
									"runningJobCount", len(periodicJob.Edges.Jobs),
									"delay", sleepTime,
								)
							}
							return nil
						}
					}
				}

				var stdErr error
				newLastScheduled = app.Clock.Now()
				if newPeriodicJobID != 0 {
					stdErr = tx.PeriodicJob.UpdateOneID(newPeriodicJobID).
						SetType(jobType).
						SetVersion(jobVersion).
						SetLastScheduledNewJob(newLastScheduled).
						Exec(ctx)
				}
				if newPeriodicJobID == 0 || ent.IsNotFound(stdErr) {
					periodicJob, stdErr := tx.PeriodicJob.Create().
						SetType(jobType).
						SetVersion(jobVersion).
						SetLastScheduledNewJob(newLastScheduled).
						Save(ctx)
					if stdErr != nil {
						return stdErr
					}
					if newPeriodicJobID == 0 {
						app.Logger.Info(
							"created new PeriodicJob as there wasn't one on startup",
							"periodicJobID", newPeriodicJobID,
							"jobType", versionedType,
						)
					} else {
						app.Logger.Error(
							"previously read PeriodicJob was not found when scheduling job, so created new periodic job",
							"periodicJobID", newPeriodicJobID,
							"jobType", versionedType,
						)
					}
					newPeriodicJobID = periodicJob.ID
				} else if stdErr != nil {
					return stdErr
				}

				_, commErr := app.Jobs.EnqueueWithModifier(
					versionedType, struct{}{},
					func(jobCreate *ent.JobCreate) *ent.JobCreate {
						return jobCreate.
							SetDue(newLastScheduled.Add(delayFunc(time.Time{}, app))).
							SetPeriodicJobID(newPeriodicJobID)
					},
					ctx,
				)
				return commErr.StandardError()
			},
		)
		if stdErr != nil {
			app.Logger.Error(
				"scheduler failed to enqueue job",
				"error", stdErr, "jobType", versionedType,
			)
		}
		periodicJobID = newPeriodicJobID
		return sleepTime
	}

	return Task{
		Init: func(engine *Engine) {
			periodicJob, stdErr := dbcommon.WithReadTx(
				context.TODO(), engine.App.Database,
				func(tx *ent.Tx, ctx context.Context) (*ent.PeriodicJob, error) {
					return tx.PeriodicJob.Query().
						Where(periodicjob.Type(jobType), periodicjob.Version(jobVersion)).
						Only(ctx)
				},
			)
			if stdErr == nil {
				lastScheduled = periodicJob.LastScheduledNewJob
				periodicJobID = periodicJob.ID
			} else if ent.IsNotFound(stdErr) {
				lastScheduled = engine.App.Clock.Now()
			} else {
				log.Fatalf("failed to read PeriodicJob object from database when starting scheduler: %v", commErr)
			}
			_ = scheduleJob(engine.App, true)
		},
		Run: func(engine *Engine) {
			extraDelay := time.Duration(0)
			for {
				select {
				case <-time.After(delayFunc(lastScheduled, engine.App) + extraDelay):
				case <-engine.RequestShutdownChan:
					return
				}
				extraDelay = scheduleJob(engine.App, false)
			}
		},
	}
}
