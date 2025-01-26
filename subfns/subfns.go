package subfns

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() *intertypes.Env {
	_ = godotenv.Load(".env")

	env := intertypes.Env{
		MOUNT_PATH:                    util.RequireEnv("MOUNT_PATH"),
		PORT:                          util.RequireIntEnv("PORT"),
		PROXY_ORIGINAL_IP_HEADER_NAME: util.RequireEnv("PROXY_ORIGINAL_IP_HEADER_NAME"),
	}
	return &env
}

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
