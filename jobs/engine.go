package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/job"
)

type Engine struct {
	App                  *common.App
	Registry             *Registry
	Running              bool
	newJobChan           chan struct{}
	waitingForJobsChan   chan struct{}
	requestShutdownChan  chan struct{}
	shutdownFinishedChan chan struct{}
	mu                   sync.Mutex
}

func NewEngine(registry *Registry) *Engine {
	return &Engine{
		App:                  registry.App,
		Registry:             registry,
		newJobChan:           make(chan struct{}, 1),
		waitingForJobsChan:   make(chan struct{}),
		requestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
	}
}

type completedJob struct {
	Object *ent.Job
	Err    *Error
}

func (engine *Engine) Listen() {
	engine.mu.Lock()
	if engine.Running {
		engine.mu.Unlock()
		panic("job engine is already running")
	}
	engine.Running = true
	engine.mu.Unlock()
	fmt.Println("job engine running")

	completedJobChan := make(chan completedJob, min(engine.App.Env.MAX_TOTAL_JOB_WEIGHT, 100))
	currentWeight := 0

	handleCompletedJob := func(completedJob completedJob) {
		currentWeight -= completedJob.Object.Weight
		if completedJob.Err == nil {
			stdErr := dbcommon.WithWriteTx(context.TODO(), engine.App.Database,
				func(tx *ent.Tx, ctx context.Context) error {
					return tx.Job.DeleteOneID(completedJob.Object.ID).Exec(ctx)
				},
			)
			if stdErr != nil {
				// A restart is unlikely to help, so we'll just have to log the error
				// TODO: log error
				fmt.Printf("failed to delete job: %v\n", stdErr.Error())
			}
			fmt.Printf("job %s completed after %d retries\n", completedJob.Object.ID, completedJob.Object.Retries)
		} else {
			jobErr := NewError(completedJob.Err)
			stdErr := dbcommon.WithWriteTx(context.TODO(), engine.App.Database,
				func(tx *ent.Tx, ctx context.Context) error {
					if completedJob.Object.Retries < len(jobErr.JobRetryBackoffs) {
						backoff := jobErr.JobRetryBackoffs[completedJob.Object.Retries]
						return tx.Job.UpdateOneID(completedJob.Object.ID).
							SetStatus("pending").
							SetDue(engine.App.Clock.Now().Add(backoff)).
							AddRetries(1).
							SetLoggedStallWarning(false).
							Exec(ctx)
					} else {
						fmt.Printf("job %s failed after %d retries. error:\n%v\n", completedJob.Object.ID, completedJob.Object.Retries, jobErr.Error())
						return tx.Job.UpdateOneID(completedJob.Object.ID).
							SetStatus("failed").
							Exec(ctx)
					}
				},
			)
			if stdErr != nil {
				// A restart is unlikely to help, so we'll just have to log the error
				// TODO: log error
				fmt.Printf("failed to mark job as failed / reset to pending: %v\n", stdErr.Error())
			}

		}
	}
	runJobs := func() bool {
		for {
			isDone := false
			shutdownSignalReceived := false
			currentJob, stdErr := dbcommon.WithReadWriteTx(
				context.TODO(), engine.App.Database,
				func(tx *ent.Tx, ctx context.Context) (*ent.Job, error) {
					currentJob, stdErr := tx.Job.Query().
						Where(job.StatusEQ("pending"), job.DueLTE(time.Now())).
						Order(ent.Asc(job.FieldStatus), ent.Desc(job.FieldPriority), ent.Asc(job.FieldDue)).
						First(context.TODO())
					if stdErr != nil {
						if ent.IsNotFound(stdErr) {
							isDone = true
							return nil, nil
						}
						return nil, ErrWrapperListen.Wrap(ErrWrapperDatabase.Wrap(stdErr)).AddCategory("query next job")
					}

					maxTotalWeightToStart := engine.App.Env.MAX_TOTAL_JOB_WEIGHT - currentWeight
					for currentWeight > maxTotalWeightToStart && currentWeight > 0 {
						select {
						case completedJob := <-completedJobChan:
							handleCompletedJob(completedJob)
						case <-engine.requestShutdownChan:
							shutdownSignalReceived = true
							return nil, nil
						}
					}
					stdErr = tx.Job.UpdateOneID(currentJob.ID).
						SetStatus("running").SetStarted(time.Now()).
						Exec(context.TODO())
					if stdErr != nil {
						return nil, ErrWrapperListen.Wrap(ErrWrapperDatabase.Wrap(stdErr)).AddCategory("update job")
					}
					return currentJob, nil
				},
			)
			if stdErr != nil {
				// This is worse than failing to update specific jobs, the program can't really continue without a job system
				fmt.Printf("job loop transaction error: %v\n", stdErr.Error())
				return true // Shutdown
			}
			if isDone {
				return false
			}
			if shutdownSignalReceived {
				return true
			}
			currentWeight += currentJob.Weight
			jobDefinition, ok := engine.Registry.jobs[common.GetVersionedType(currentJob.Type, currentJob.Version)]
			if ok {
				if jobDefinition.NoParallelize {
					// TODO: this doesn't work because the engine's transaction doesn't commit until this has finished
					engine.runJob(jobDefinition, currentJob, completedJobChan)
				} else {
					go engine.runJob(jobDefinition, currentJob, completedJobChan)
				}
			} else { // Note: this shouldn't happen
				completedJobChan <- completedJob{
					Object: currentJob,
					Err:    NewError(ErrWrapperRunJob.Wrap(ErrUnknownJobType)),
				}
			}
			// Otherwise continue processing jobs
		}
	}

listenLoop:
	for {
		fmt.Println("run job loop")
		if runJobs() { // Shutdown
			break listenLoop
		}
		// TODO: check for stalled jobs, do a similar thing to scheduled jobs. Wait until they should have finished
		// Maybe schedule a job to check for stalled jobs every 5 minutes or so, because checking after all the jobs have finished running might be too late if there are a lot

		maxWaitTime := engine.App.Env.JOB_POLL_INTERVAL
		nextJob, stdErr := dbcommon.WithReadTx(context.TODO(), engine.App.Database,
			func(tx *ent.Tx, ctx context.Context) (*ent.Job, error) {
				return tx.Job.Query().
					Where(job.StatusEQ("pending")).
					Order(ent.Asc(job.FieldDue)).
					First(ctx)
			})
		if stdErr == nil {
			timeUntil := time.Until(nextJob.Due)
			if timeUntil < maxWaitTime {
				maxWaitTime = timeUntil
			}
		} else if !ent.IsNotFound(stdErr) {
			fmt.Printf("failed to query next due job: %v\n", stdErr.Error())
			break listenLoop
		}

		for currentWeight > 0 {
			select {
			case completedJob := <-completedJobChan:
				handleCompletedJob(completedJob)
			case <-engine.requestShutdownChan:
				break listenLoop
			}
		}

		close(engine.waitingForJobsChan)
		engine.mu.Lock()
		engine.waitingForJobsChan = make(chan struct{})
		engine.mu.Unlock()

		fmt.Println("waiting for new jobs")
		select {
		case <-time.After(maxWaitTime):
		case <-engine.newJobChan:
		case <-engine.requestShutdownChan:
			break listenLoop
		}
	}
	for currentWeight > 0 {
		completedJob := <-completedJobChan
		handleCompletedJob(completedJob)
	}
	close(engine.requestShutdownChan)
	close(engine.shutdownFinishedChan)
}
func (engine *Engine) runJob(
	jobDefinition *Definition, job *ent.Job,
	completedJobChan chan completedJob,
) {
	stdErr := jobDefinition.Handler(&Context{
		Definition: jobDefinition,
		Context:    context.TODO(),
		Body:       []byte(job.Data),
	})
	completedJobChan <- completedJob{
		Object: job,
		Err:    NewError(ErrWrapperRunJob.Wrap(stdErr)),
	}
}

func (engine *Engine) WaitForJobs() {
	<-engine.waitingForJobsChan
}

func (engine *Engine) Shutdown() {
	// TODO: timeout?
	// TODO: what if it's not running?
	fmt.Println("job engine shutting down")
	// TODO: this panics if the engine has to shut itself down due to an error
	engine.requestShutdownChan <- struct{}{}
	fmt.Println("job engine finishing jobs")
	<-engine.shutdownFinishedChan
	fmt.Println("job engine stopped")
}

func (engine *Engine) Enqueue(
	versionedType string,
	data any,
	ctx context.Context,
) (uuid.UUID, *common.Error) {
	_, ok := engine.Registry.jobs[versionedType]
	if !ok {
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(ErrUnknownJobType)
	}
	encoded, commErr := engine.Registry.Encode(
		versionedType,
		data,
	)
	if commErr != nil {
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(commErr)
	}

	return engine.EnqueueEncoded(versionedType, encoded, ctx)
}

func (engine *Engine) EnqueueEncoded(
	versionedType string,
	encodedData string,
	ctx context.Context,
) (uuid.UUID, *common.Error) {
	jobDefinition, ok := engine.Registry.jobs[versionedType]
	if !ok {
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(ErrUnknownJobType)
	}

	jobType, version, commErr := common.ParseVersionedType(versionedType)
	if commErr != nil { // This shouldn't happen because of the Encode call but just in case
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(commErr)
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(ErrNoTxInContext)
	}
	job, stdErr := tx.Job.Create().
		SetType(jobType).
		SetVersion(version).
		SetPriority(jobDefinition.Priority).
		SetWeight(jobDefinition.Weight).
		SetData(encodedData).
		Save(ctx)
	if stdErr != nil {
		return uuid.UUID{}, ErrWrapperEnqueue.Wrap(ErrWrapperDatabase.Wrap(stdErr))
	}

	// Otherwise the engine will look for the job before it's committed
	tx.OnCommit(func(next ent.Committer) ent.Committer {
		return ent.CommitFunc(
			func(ctx context.Context, tx *ent.Tx) error {
				stdErr := next.Commit(ctx, tx)
				if stdErr != nil {
					return stdErr
				}

				fmt.Println("sending new job signal")
				select {
				case engine.newJobChan <- struct{}{}:
				default:
					fmt.Println("new job signal already sent")
				}
				return nil
			},
		)
	})

	return job.ID, nil
}
