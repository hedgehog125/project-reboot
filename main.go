package main

import (
	"github.com/hedgehog125/project-reboot/subfns"
	"github.com/jonboulle/clockwork"
)

func main() {
	env := subfns.LoadEnvironmentVariables()
	clock := clockwork.NewRealClock()

	state := subfns.InitState()
	dbClient := subfns.OpenDatabase(env)
	messengerGroup := subfns.ConfigureMessengers(env)
	engine := subfns.ConfigureServer(state, dbClient, messengerGroup, clock, env)
	scheduler := subfns.ConfigureScheduler(clock, state)

	subfns.RunScheduler(scheduler)
	server := subfns.RunServer(engine, env)

	subfns.ConfigureShutdown(
		subfns.NewShutdownTask(func() {
			subfns.ShutdownScheduler(scheduler)
		}, true),
		subfns.NewShutdownTask(func() {
			subfns.ShutdownServer(server)
		}, true),
		subfns.NewShutdownTask(func() {
			dbClient.Close()
		}, false),
	)
}
