package services

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
)

// TODO: replace with job scheduler. Maybe the state struct should be removed and the admin token rotation becomes a simple loop?

func NewScheduler(app *common.App) common.SchedulerService {
	scheduler, err := gocron.NewScheduler(
		gocron.WithClock(app.Clock),
		gocron.WithLocation(time.UTC),
		gocron.WithStopTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("couldn't start scheduler. error:\n %v", err.Error())
	}

	addJobs(scheduler, app)
	return &schedulerService{
		scheduler: scheduler,
	}
}
func addJobs(scheduler gocron.Scheduler, app *common.App) {
	// Once an hour
	mustAddJob(scheduler, gocron.CronJob("0 * * * *", false), gocron.NewTask(core.UpdateAdminCode, app.State))
}

type schedulerService struct {
	scheduler gocron.Scheduler
}

func (service *schedulerService) Start() {
	service.scheduler.Start()
}

func (service *schedulerService) Shutdown() {
	err := service.scheduler.Shutdown()
	if err != nil {
		// TODO: replace with slog when the new scheduler is written
		fmt.Printf("warning: an error occurred while shutting down the scheduler:\n%v\n", err.Error())
	}
}

func mustAddJob(scheduler gocron.Scheduler, jobDefinition gocron.JobDefinition, task gocron.Task, jobOptions ...gocron.JobOption) gocron.Job {
	job, err := scheduler.NewJob(jobDefinition, task, jobOptions...)
	if err != nil {
		log.Fatalf("couldn't create job. error: %v", err.Error())
	}

	return job
}
