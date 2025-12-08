package services

import (
	"context"
	"runtime"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/schedulers"
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
					return app.Core.SendActiveSessionReminders(ctx)
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
					wrappedErr := app.Core.DeleteExpiredSessions(ctx)
					if wrappedErr != nil {
						return wrappedErr
					}
					wrappedErr = app.TwoFactorActions.DeleteExpiredActions(ctx)
					if wrappedErr != nil {
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
	if app.Env.FULL_GC_INTERVAL > 0 {
		engine.Register(
			func(taskContext *schedulers.TaskContext) {
				app.Logger.Info("running full garbage collection...")
				runtime.GC()
			},
			schedulers.SimpleFixedInterval(app.Env.FULL_GC_INTERVAL),
		)
	}

	return &Scheduler{
		Engine: engine,
	}
}
func (scheduler *Scheduler) Start() {
	go scheduler.Engine.Run()
}
