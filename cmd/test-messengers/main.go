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
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
	"github.com/hedgehog125/project-reboot/services"
)

func main() {
	username := flag.String("username", "", "the name of the user to send messages to")
	flag.Parse()
	if *username == "" {
		log.Fatalf("missing required argument \"username\"")
	}

	env := services.LoadEnvironmentVariables()
	databaseService := services.NewDatabase(env)
	databaseService.Start()
	defer databaseService.Shutdown()

	discord := messengers.NewDiscord(env)
	var userInfo *common.UserContacts
	err := dbcommon.WithTx(
		context.Background(), databaseService,
		func(tx *ent.Tx) error {
			var err error
			userInfo, err = messengerscommon.ReadUserContacts(*username, ent.NewTxContext(context.Background(), tx))
			return err
		},
	)
	if err != nil {
		panic(fmt.Sprintf("couldn't read user. error:\n%v", err.Error()))
	}
	err = discord.Send(common.Message{
		Type: common.MessageTest,
		User: userInfo,
	})
	if err != nil {
		panic(fmt.Sprintf("couldn't send Discord message. error:\n%v", err.Error()))
	}
	fmt.Fprintln(os.Stdout, "sent Discord message")
}
