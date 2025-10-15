package services

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/schedulers"
)

type Scheduler struct {
	*schedulers.Engine
}

func NewScheduler(app *common.App) *Scheduler {
	engine := schedulers.NewEngine(app)

	engine.Register(func(taskContext *schedulers.TaskContext) {
		app.Core.RotateAdminCode()
	}, schedulers.SimpleFixedInterval(1*time.Hour))
	engine.Register(
		func(taskContext *schedulers.TaskContext) {
			app.Logger.Info("sending active session reminders...")
			stdErr := dbcommon.WithWriteTx(
				taskContext.Context, taskContext.App.Database,
				func(tx *ent.Tx, ctx context.Context) error {
					return app.Core.SendActiveSessionReminders(ctx).StandardError()
				},
			)
			if stdErr != nil {
				taskContext.App.Logger.Error(
					"unable to send active session reminders!",
					"error", stdErr,
				)
			}
		},
		schedulers.PersistentFixedInterval(
			"SEND_ACTIVE_SESSION_REMINDERS",
			app.Env.ACTIVE_SESSION_REMINDER_INTERVAL,
		),
	)
	engine.Register(
		func(taskContext *schedulers.TaskContext) {
			app.Logger.Info("running cleanup task...")
			stdErr := dbcommon.WithWriteTx(
				taskContext.Context, taskContext.App.Database,
				func(tx *ent.Tx, ctx context.Context) error {
					commErr := app.Core.DeleteExpiredSessions(ctx)
					if commErr != nil {
						return commErr
					}
					commErr = app.TwoFactorActions.DeleteExpiredActions(ctx)
					if commErr != nil {
						return nil
					}

					return nil
				},
			)
			if stdErr != nil {
				taskContext.App.Logger.Error(
					"cleanup task's transaction failed",
					"error", stdErr,
				)
			}
			app.RateLimiter.DeleteInactiveUsers()
		},
		schedulers.SimpleFixedInterval(app.Env.CLEAN_UP_INTERVAL),
	)

	return &Scheduler{
		Engine: engine,
	}
}
func (scheduler *Scheduler) Start() {
	go scheduler.Engine.Run()
}
