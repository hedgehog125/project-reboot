package main

import (
	"flag"
	"fmt"
	"log"

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
	dbClient := services.OpenDatabase(env)
	defer dbClient.Close()

	discord := messengers.NewDiscord(env)
	userInfo, err := messengers.ReadUserInfo(*username, dbClient)
	if err != nil {
		log.Fatalf("couldn't read user. error:\n%v", err.Error())
	}
	err = discord.Send(messengers.Message{
		Type: messengers.MessageTest,
		User: userInfo,
	})
	if err != nil {
		log.Fatalf("couldn't send Discord message. error:\n%v", err.Error())
	}
	fmt.Println("sent Discord message")
}
