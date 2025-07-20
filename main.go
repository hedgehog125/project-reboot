package main

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
)

func main() {
	app := common.App{
		Env:   services.LoadEnvironmentVariables(),
		Clock: clockwork.NewRealClock(),
	}

	app.State = services.InitState()
	app.Database = services.NewDatabase(app.Env)
	app.Database.Start()
	app.Messenger = services.NewMessenger(app.Env)
	app.Scheduler = services.NewScheduler(&app)
	app.Jobs = services.NewJob(&app)
	app.Server = services.NewServer(&app)

	// TODO: add panic handling
	app.Scheduler.Start()
	app.Server.Start()
	app.Jobs.Start()

	services.ConfigureShutdown(
		services.NewShutdownTask(func() {
			app.Scheduler.Shutdown()
		}, true),
		services.NewShutdownTask(func() {
			app.Server.Shutdown()
		}, true),

		services.NewShutdownTask(func() {
			app.Jobs.Shutdown()
		}, false),
		services.NewShutdownTask(func() {
			app.Database.Shutdown()
		}, false),
	)
}
