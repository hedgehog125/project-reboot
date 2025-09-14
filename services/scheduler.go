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

	engine.Register(func(app *common.App) {
		core.UpdateAdminCode(app.State)
	}, schedulers.FixedInterval(1*time.Hour))

	return &Scheduler{
		Engine: engine,
	}
}
func (scheduler *Scheduler) Start() {
	go scheduler.Engine.Run()
}
