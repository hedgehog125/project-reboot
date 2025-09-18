package services

import (
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/schedulers"
)

type Scheduler struct {
	*schedulers.Engine
}

func NewScheduler(app *common.App) *Scheduler {
	engine := schedulers.NewEngine(app)

	engine.Register(func(taskContext *schedulers.TaskContext) {
		core.UpdateAdminCode(app.State)
	}, schedulers.SimpleFixedInterval(1*time.Hour))
	engine.Register(func(taskContext *schedulers.TaskContext) {
		app.Logger.Info("it's a new day!")
	}, schedulers.PersistentFixedInterval("TEST_PERIODIC_JOB", 24*time.Hour))

	return &Scheduler{
		Engine: engine,
	}
}
func (scheduler *Scheduler) Start() {
	go scheduler.Engine.Run()
}
