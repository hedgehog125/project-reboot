package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
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
	userInfo, err := common.ReadMessageUserInfo(*username, databaseService.Client())
	if err != nil {
		log.Fatalf("couldn't read user. error:\n%v", err.Error())
	}
	err = discord.Send(common.Message{
		Type: common.MessageTest,
		User: userInfo,
	})
	if err != nil {
		log.Fatalf("couldn't send Discord message. error:\n%v", err.Error())
	}
	fmt.Println("sent Discord message")
}
