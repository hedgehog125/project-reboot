package schedulers

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	MaxTaskRunTime   = 10 * time.Second
	DelayFuncTimeout = 350 * time.Millisecond // How long it has to calculate the next delay
)

// TODO: this abstraction seems unnecessary
type Task struct {
	Init func(engine *Engine)
	Run  func(engine *Engine)
}
type TaskCallback func(taskContext *TaskContext) // TODO: why does not having an equals here work? Why does it now mean you get the proper function snippet?
type TaskContext struct {
	App      *common.App
	CallTime time.Time
	Context  context.Context
}

func NewTask(callback TaskCallback, delayFunc DelayFunc) Task {
	var commit CommitDelayFunc
	nextRun := time.Time{}
	return Task{
		Init: func(engine *Engine) {
			ctx, cancel := context.WithTimeout(context.Background(), DelayFuncTimeout)
			defer cancel()
			nextRun, commit = delayFunc(&DelayFuncContext{
				App:     engine.App,
				LastRan: time.Time{},
				Context: ctx,
			})
		},
		Run: func(engine *Engine) {
			tick := func() {
				ctx, cancel := context.WithTimeout(context.Background(), MaxTaskRunTime)
				defer cancel()
				callback(&TaskContext{
					App:      engine.App,
					CallTime: nextRun,
					Context:  ctx,
				})
				cancel()

				ctx, cancel = context.WithTimeout(context.Background(), DelayFuncTimeout)
				defer cancel()
				commit(nextRun, ctx)

				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), DelayFuncTimeout)
				defer cancel()
				nextRun, commit = delayFunc(&DelayFuncContext{
					App:     engine.App,
					LastRan: nextRun,
					Context: ctx,
				})
			}

			for {
				select {
				case <-engine.App.Clock.After(engine.App.Clock.Until(nextRun)):
				case <-engine.RequestShutdownChan:
					return
				}
				tick()
			}
		},
	}
}

// Note: if you need more advanced behaviour, like scheduling multiple jobs or providing bodies, use NewTask instead
func NewJobTask(versionedType string, delayFunc DelayFunc) Task {
	return NewTask(func(taskContext *TaskContext) {
		_, commErr := taskContext.App.Jobs.Enqueue(
			versionedType,
			struct{}{},
			taskContext.Context,
		)
		if commErr != nil {
			taskContext.App.Logger.Error(
				"failed to enqueue periodic job",
				"error", commErr,
				"jobType", versionedType,
			)
		}
	}, delayFunc)
}
