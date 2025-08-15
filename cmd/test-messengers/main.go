package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/services"
)

func main() {
	username := flag.String("username", "", "the name of the user to send messages to")
	flag.Parse()
	if *username == "" {
		log.Fatalf("missing required argument \"username\"")
	}

	app := &common.App{
		Env: services.LoadEnvironmentVariables(),
	}
	app.Database = services.NewDatabase(app.Env)
	app.Database.Start()
	defer app.Database.Shutdown()

	{
		messengerService := services.NewMessengers(app)
		app.Messengers = messengerService
		app.Jobs = services.NewJobs(app, messengerService.RegisterJobs)
		defer app.Jobs.Shutdown()
	}

	userOb, stdErr := dbcommon.WithReadTx(
		context.Background(), app.Database,
		func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
			userOb, stdErr := tx.User.Query().
				Where(user.Username(*username)).
				Only(ctx)
			return userOb, stdErr
		},
	)
	if stdErr != nil {
		panic(fmt.Sprintf("couldn't read user. error:\n%v", stdErr.Error()))
	}
	_, commErr := app.Messengers.SendUsingAll(&common.Message{
		Type: common.MessageTest,
		User: userOb,
	}, context.Background())
	if commErr != nil {
		panic(fmt.Sprintf("couldn't send queue message. error:\n%v", commErr.Error()))
	}
	fmt.Fprintln(os.Stdout, "waiting for message jobs to run...")
	app.Jobs.WaitForJobs()
}
