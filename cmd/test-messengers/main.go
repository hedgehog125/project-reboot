package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/subfns"
)

func main() {
	username := flag.String("username", "", "the name of the user to send messages to")
	flag.Parse()
	if *username == "" {
		log.Fatalf("missing required argument \"username\"")
	}

	env := subfns.LoadEnvironmentVariables()
	dbClient := subfns.OpenDatabase(env)
	defer dbClient.Close()

	discord := messengers.NewDiscord(env)
	row, err := dbClient.User.Query().
		Where(user.Username(*username)).
		Select(user.FieldAlertDiscordId).
		Only(context.Background())
	if err != nil {
		log.Fatalf("discord error: %v", err.Error())
	}
	err = discord.SendBatch([]messengers.Message{
		{
			Type: messengers.MessageTest,
			User: row,
		},
	})
	if err != nil {
		log.Fatalf("couldn't send Discord message. error:\n%v", err.Error())
	}
	fmt.Println("sent Discord message")
}
