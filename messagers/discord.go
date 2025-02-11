package messagers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/intertypes"
)

type discord struct {
	env     *intertypes.Env
	session *discordgo.Session
}

type Discord interface {
	Messager
}

func NewDiscord(env *intertypes.Env) Discord {
	session, err := discordgo.New("Bot " + env.DISCORD_TOKEN)
	if err != nil {
		log.Fatalf("error creating Discord session:\n%v", err)
	}
	return &discord{
		env:     env,
		session: session,
	}
}

func (discord *discord) SendBatch(messages []Message) error {
	err := discord.session.Open()
	if err != nil {
		return err
	}
	defer discord.session.Close()

	for _, message := range messages {
		preparedMessage, err := prepareMessage(message)
		if err != nil {
			return err
		}

		channel, err := discord.session.UserChannelCreate(message.User.AlertDiscordId)
		if err != nil {
			return err
		}
		_, err = discord.session.ChannelMessageSend(channel.ID, preparedMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func prepareMessage(message Message) (string, error) {
	switch message.Type {
	case MessageLogin:
		return fmt.Sprintf("Login attempt"), nil
	}
	return "", fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type)
}
