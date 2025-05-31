package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hedgehog125/project-reboot/common"
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
	defer databaseService.Shutdown()

	discord := messengers.NewDiscord(env)
	userInfo, err := messengerscommon.ReadMessageUserInfo(*username, databaseService.Client())
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
