package messengers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/common"
)

type discord struct {
	env     *common.Env
	session *discordgo.Session
}

type Discord interface {
	Messenger
}

func NewDiscord(env *common.Env) Discord {
	session, err := discordgo.New("Bot " + env.DISCORD_TOKEN)
	if err != nil {
		log.Fatalf("error creating Discord session:\n%v", err)
	}
	return &discord{
		env:     env,
		session: session,
	}
}

func (discord *discord) Id() string {
	return "discord"
}

func (discord *discord) Send(message common.Message) error {
	err := discord.session.Open()
	if err != nil {
		return err
	}
	defer discord.session.Close()

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

	return nil
}

func prepareMessage(message common.Message) (string, error) {
	switch message.Type {
	case common.MessageLogin:
		return fmt.Sprintf("Login attempt"), nil
	case common.MessageTest:
		return "Test message", nil
	}
	return "", fmt.Errorf("message type \"%v\" hasn't been implemented", message.Type)
}
