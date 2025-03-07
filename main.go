package main

import (
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
)

func main() {
	env := services.LoadEnvironmentVariables()
	clock := clockwork.NewRealClock()

	state := services.InitState()
	dbClient := services.OpenDatabase(env)
	messengerGroup := services.ConfigureMessengers(env)
	engine := services.ConfigureServer(state, dbClient, messengerGroup, clock, env)
	scheduler := services.ConfigureScheduler(clock, state)

	services.RunScheduler(scheduler)
	server := services.RunServer(engine, env)

	services.ConfigureShutdown(
		services.NewShutdownTask(func() {
			services.ShutdownScheduler(scheduler)
		}, true),
		services.NewShutdownTask(func() {
			services.ShutdownServer(server)
		}, true),
		services.NewShutdownTask(func() {
			dbClient.Close()
		}, false),
	)
}
