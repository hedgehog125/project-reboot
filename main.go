package main

import (
	"github.com/hedgehog125/project-reboot/subfns"
	"github.com/jonboulle/clockwork"
)

func main() {
	env := subfns.LoadEnvironmentVariables()
	clock := clockwork.NewRealClock()

	dbClient := subfns.OpenDatabase(env)
	state := subfns.InitState()
	engine := subfns.ConfigureServer(state, dbClient, clock, env)
	scheduler := subfns.ConfigureScheduler(clock, state)

	subfns.RunScheduler(scheduler)
	server := subfns.RunServer(engine, env)

	// { // TODO: remove
	// 	discord := messagers.NewDiscord(env)
	// 	row, err := dbClient.User.Query().
	// 		Where(user.Username("user1")).
	// 		Select(user.FieldAlertDiscordId).
	// 		Only(context.Background())
	// 	if err != nil {
	// 		log.Fatalf("discord error: %v", err.Error())
	// 	}
	// 	fmt.Printf("%v", discord.SendBatch([]messagers.Message{
	// 		{
	// 			Type: messagers.MessageLogin,
	// 			User: row,
	// 		},
	// 	}))
	// }

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
