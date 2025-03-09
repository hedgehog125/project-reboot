package services

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

func ConfigureShutdown(tasks ...*ShutdownTask) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nshutting down...")
	var wg sync.WaitGroup
	for _, task := range tasks {
		if task.Concurrent {
			wg.Add(1)
			go func() {
				task.Callback()
				wg.Done()
			}()
		} else {
			wg.Wait()
			task.Callback()
		}
	}
	wg.Wait()
	fmt.Println("shut down.")
}
