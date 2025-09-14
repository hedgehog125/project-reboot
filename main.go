package main

import (
	"log/slog"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
)

func main() {
	app := &common.App{
		Env:   services.LoadEnvironmentVariables(),
		Clock: clockwork.NewRealClock(),
	}
	shutdownService := services.NewShutdown(
		app,
		services.NewShutdownTask(func() {
			if app.Scheduler != nil {
				app.Scheduler.Shutdown()
			}
		}, true),
		services.NewShutdownTask(func() {
			if app.Server != nil {
				app.Server.Shutdown()
			}
		}, true),

		services.NewShutdownTask(func() {
			if app.Jobs != nil {
				app.Jobs.Shutdown()
			}
		}, false),
		services.NewShutdownTask(func() {
			if app.Logger != nil {
				app.Logger.Shutdown() // Note: logs will still be written to the console after this, just not stored
			}
		}, false),
		services.NewShutdownTask(func() {
			if app.Database != nil {
				app.Database.Shutdown()
			}
		}, false),
	)
	app.ShutdownService = shutdownService

	app.State = services.InitState()
	{
		logger := services.NewLogger(app)
		app.Logger = logger
		slog.SetDefault(logger.Logger)
	}
	app.Database = services.NewDatabase(app)
	app.Database.Start()
	app.Logger.Start()
	app.Scheduler = services.NewScheduler(app)
	app.TwoFactorActions = services.NewTwoFactorActions(app)
	{
		messengerService := services.NewMessengers(app)
		app.Messengers = messengerService
		app.Jobs = services.NewJobs(app, messengerService.RegisterJobs)
	}
	app.Server = services.NewServer(app)

	app.Scheduler.Start()
	app.Server.Start()
	app.Jobs.Start()

	shutdownService.Listen()
}
