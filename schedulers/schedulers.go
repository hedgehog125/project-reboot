package schedulers

import (
	"sync"

	"github.com/hedgehog125/project-reboot/common"
)

type Engine struct {
	App   *common.App
	tasks []Task
	// Note: outside of the schedulers package, this should only be used by LoopFuncs to listen for shutdown requests
	RequestShutdownChan  chan struct{}
	shutdownFinishedChan chan struct{}
	shutdownWg           sync.WaitGroup
}

func NewEngine(app *common.App) *Engine {
	return &Engine{
		App:                  app,
		tasks:                make([]Task, 0),
		RequestShutdownChan:  make(chan struct{}),
		shutdownFinishedChan: make(chan struct{}),
	}
}
func (engine *Engine) Register(fn func(app *common.App), delayFunc DelayFunc) {
	engine.RegisterTask(NewTask(fn, delayFunc))
}
func (engine *Engine) RegisterJob(versionedName string, delayFunc DelayFunc, maxConcurrentRuns int) {
	engine.RegisterTask(NewJobTask(versionedName, delayFunc, maxConcurrentRuns))
}
func (engine *Engine) RegisterTask(task Task) {
	engine.tasks = append(engine.tasks, task)
}
func (engine *Engine) Run() {
	for _, task := range engine.tasks {
		if task.Init != nil {
			task.Init(engine)
		}
		engine.shutdownWg.Go(func() {
			task.Run(engine)
		})
	}
	engine.shutdownWg.Wait()
	close(engine.shutdownFinishedChan)
}
func (engine *Engine) Shutdown() {
	engine.App.Logger.Info("scheduler shutting down...")
	close(engine.RequestShutdownChan)
	<-engine.shutdownFinishedChan
	engine.App.Logger.Info("scheduler stopped")
}
