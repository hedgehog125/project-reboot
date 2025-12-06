package services

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hedgehog125/project-reboot/common"
)

type ShutdownTask struct {
	Callback   func()
	Concurrent bool
}

func NewShutdownTask(callback func(), concurrent bool) *ShutdownTask {
	return &ShutdownTask{
		Callback:   callback,
		Concurrent: concurrent,
	}
}

type Shutdown struct {
	app                   *common.App
	tasks                 []*ShutdownTask
	shutdownStartedChan   chan struct{}
	shutdownCompletedChan chan struct{}
	shutdownOnce          sync.Once
}

func NewShutdown(app *common.App, tasks ...*ShutdownTask) *Shutdown {
	return &Shutdown{
		app:                   app,
		tasks:                 tasks,
		shutdownStartedChan:   make(chan struct{}),
		shutdownCompletedChan: make(chan struct{}),
	}
}
func (shutdown *Shutdown) Listen() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigChan:
		shutdown.Shutdown("")
	case <-shutdown.shutdownStartedChan:
		<-shutdown.shutdownCompletedChan
	}
}

// Unlike Listen, this doesn't listen for a SIGINT. Mostly used for tests
func (shutdown *Shutdown) ListenForShutdownCall() {
	<-shutdown.shutdownStartedChan
	<-shutdown.shutdownCompletedChan
}
func (shutdown *Shutdown) Shutdown(reason string) {
	shutdown.shutdownOnce.Do(func() {
		close(shutdown.shutdownStartedChan)
		if reason == "" {
			shutdown.app.Logger.Info("shutting down...")
		} else {
			shutdown.app.Logger.Error("shutting down due to a critical error!", "reason", reason)
		}
		var wg sync.WaitGroup
		for _, task := range shutdown.tasks {
			if task.Concurrent {
				wg.Go(func() {
					task.Callback()
				})
			} else {
				wg.Wait()
				task.Callback()
			}
		}
		wg.Wait()
		shutdown.app.Logger.Info("shut down.")
		close(shutdown.shutdownCompletedChan)
		if reason != "" {
			os.Exit(1)
		}
	})
}
