package main

import (
	"github.com/hedgehog125/project-reboot/subfns"
	"github.com/jonboulle/clockwork"
)

func main() {
	env := subfns.LoadEnvironmentVariables()
	clock := clockwork.NewRealClock()

	_ = subfns.OpenDatabase(env)
	engine := subfns.ConfigureServer(env)
	state := subfns.InitState()
	scheduler := subfns.ConfigureScheduler(clock, state)

	subfns.RunScheduler(scheduler)

	go subfns.RunServer(engine, env)

	subfns.ConfigureShutdown(
		subfns.NewShutdownTask(func() {
			subfns.ShutdownScheduler(scheduler)
		}, true),
	)
}
