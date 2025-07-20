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
	"github.com/hedgehog125/project-reboot/jobs/jobscommon"
)

type Engine struct {
	App                  *common.App
	Registry             *Registry
	Running              bool
	newJobChan           chan struct{}
	requestShutdownChan  chan struct{}
	shutdownFinishedChan chan struct{}
	mu                   sync.Mutex
}

func NewEngine(registry *Registry) *Engine {
	return &Engine{
		App:                  registry.App,
		Registry:             registry,
		newJobChan:           make(chan struct{}, 1),
		requestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
	}
}

type completedJob struct {
	ID     uuid.UUID
	Weight int
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

	dbClient := engine.App.Database.Client()
	completedJobChan := make(chan completedJob, min(engine.App.Env.MAX_TOTAL_JOB_WEIGHT, 100))
	currentWeight := 0

	handleCompletedJob := func(completedJob completedJob) {
		currentWeight -= completedJob.Weight
		stdErr := dbClient.Job.DeleteOneID(completedJob.ID).Exec(context.TODO())
		if stdErr != nil {
			// TODO: log error
			panic(fmt.Sprintf("failed to delete job: %v\n", stdErr.Error()))
		}
	}
	runJobs := func() bool {
		for {
			isDone := false
			shutdownSignalReceived := false
			stdErr := dbcommon.WithTx(
				context.TODO(), engine.App.Database,
				func(tx *ent.Tx) error {
					currentJob, stdErr := tx.Job.Query().
						Where(job.StatusEQ("pending"), job.DueLTE(time.Now())).
						Order(ent.Asc(job.FieldStatus), ent.Desc(job.FieldPriority), ent.Asc(job.FieldDue)).
						First(context.TODO())
					if stdErr != nil {
						if ent.IsNotFound(stdErr) {
							isDone = true
							return nil
						}
						return ErrWrapperDatabase.Wrap(stdErr).AddCategory("query next job")
					}

					maxTotalWeightToStart := engine.App.Env.MAX_TOTAL_JOB_WEIGHT - currentWeight
					for currentWeight > maxTotalWeightToStart && currentWeight > 0 {
						select {
						case completedJob := <-completedJobChan:
							handleCompletedJob(completedJob)
						case <-engine.requestShutdownChan:
							shutdownSignalReceived = true
							return nil
						}
					}
					stdErr = tx.Job.UpdateOneID(currentJob.ID).
						SetStatus("running").SetStarted(time.Now()).
						Exec(context.TODO())
					if stdErr != nil {
						return ErrWrapperDatabase.Wrap(stdErr).AddCategory("update job")
					}
					currentWeight += currentJob.Weight
					go engine.runJob(currentJob, completedJobChan)
					return nil
				},
			)
			if stdErr != nil {
				// TODO: retry up to 3 times, logging error each time, then just crash
				fmt.Println(stdErr.Error())
				return true // Shutdown
			}
			if isDone {
				return false
			}
			if shutdownSignalReceived {
				return true
			}
			// Otherwise continue processing jobs
		}
	}

	// TODO: handle panics
listenLoop:
	for {
		if runJobs() { // Shutdown
			break listenLoop
		}
		// TODO: check for stalled jobs, do a similar thing to scheduled jobs. Wait until they should have finished
		// Maybe schedule a job to check for stalled jobs every 5 minutes or so

		maxWaitTime := engine.App.Env.JOB_POLL_INTERVAL
		nextJob, stdErr := dbClient.Job.Query().
			Where(job.StatusEQ("pending")).
			Order(ent.Asc(job.FieldDue)).
			First(context.TODO())
		if stdErr == nil {
			timeUntil := time.Until(nextJob.Due)
			if timeUntil < maxWaitTime {
				maxWaitTime = timeUntil
			}
		} else if !ent.IsNotFound(stdErr) {
			// TODO: retry up to 3 times, logging error each time, then just crash
			fmt.Printf("failed to query next due job: %v\n", stdErr.Error())
			break listenLoop
		}

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
	close(engine.shutdownFinishedChan)
}
func (engine *Engine) runJob(job *ent.Job, completedJobChan chan completedJob) {
	jobDefinition, ok := engine.Registry.jobs[jobscommon.GetVersionedType(job.Type, job.Version)]
	if !ok { // Note: this shouldn't happen
		completedJobChan <- completedJob{
			ID:     job.ID,
			Weight: job.Weight,
			Err:    NewError(ErrUnknownJobType.AddCategory(ErrTypeRunJob)),
		}
		return
	}

	stdErr := jobDefinition.Handler(&Context{
		Definition: jobDefinition,
		Context:    context.TODO(),
		Body:       []byte(job.Data),
	})
	completedJobChan <- completedJob{
		ID:     job.ID,
		Weight: job.Weight,
		Err:    NewError(stdErr).AddCategory(ErrTypeRunJob), // TODO: are these categories correct?
	}
}

func (engine *Engine) Shutdown() {
	// TODO: timeout?
	// TODO: what if it's not running?
	engine.requestShutdownChan <- struct{}{}
	<-engine.shutdownFinishedChan
}

func (engine *Engine) Enqueue(
	versionedType string,
	data any,
	ctx context.Context,
) (uuid.UUID, *common.Error) {
	jobDefinition, ok := engine.Registry.jobs[versionedType]
	if !ok {
		return uuid.UUID{}, ErrUnknownJobType.AddCategory(ErrTypeEnqueue)
	}
	encoded, commErr := engine.Registry.Encode(
		versionedType,
		data,
	)
	if commErr != nil {
		return uuid.UUID{}, commErr.AddCategory(ErrTypeEnqueue)
	}

	jobType, version, commErr := jobscommon.ParseVersionedType(versionedType)
	if commErr != nil { // This shouldn't happen because of the Encode call but just in case
		return uuid.UUID{}, commErr.AddCategory(ErrTypeEnqueue)
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return uuid.UUID{}, ErrNoTxInContext.AddCategory(ErrTypeEnqueue)
	}
	action, stdErr := tx.Job.Create().
		SetType(jobType).
		SetVersion(version).
		SetPriority(jobDefinition.Priority).
		SetWeight(jobDefinition.Weight).
		SetData(encoded).
		Save(ctx)
	if stdErr != nil {
		return uuid.UUID{}, ErrWrapperDatabase.Wrap(stdErr).AddCategory(ErrTypeEnqueue)
	}

	select {
	case engine.newJobChan <- struct{}{}:
	default:
	}

	return action.ID, nil
}
