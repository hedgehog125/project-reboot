package subfns

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/jonboulle/clockwork"
)

func ConfigureScheduler(clock clockwork.Clock, state intertypes.State) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler(
		gocron.WithClock(clock),
		gocron.WithLocation(time.UTC),
		gocron.WithStopTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("couldn't start scheduler. error:\n %v", err.Error())
	}

	addJobs(scheduler, state)
	return scheduler
}
func addJobs(scheduler gocron.Scheduler, state intertypes.State) {
	// Once an hour
	mustAddJob(scheduler, gocron.CronJob("0 * * * *", false), gocron.NewTask(core.UpdateAdminCode, state))

	mustAddJob(scheduler, gocron.CronJob("* * * * *", false), gocron.NewTask(func() {
		fmt.Println("waiting")
		time.Sleep(150 * time.Second)
		fmt.Println("done")
	})) // Every minute
}

func RunScheduler(scheduler gocron.Scheduler) {
	scheduler.Start()
}

func ShutdownScheduler(scheduler gocron.Scheduler) {
	err := scheduler.StopJobs()
	if err == nil {
		err = scheduler.Shutdown()
	}
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the scheduler:\n%v", err.Error())
	}
}

func mustAddJob(scheduler gocron.Scheduler, jobDefinition gocron.JobDefinition, task gocron.Task, jobOptions ...gocron.JobOption) gocron.Job {
	job, err := scheduler.NewJob(jobDefinition, task, jobOptions...)
	if err != nil {
		log.Fatalf("couldn't create job. error: %v", err.Error())
	}

	return job
}
