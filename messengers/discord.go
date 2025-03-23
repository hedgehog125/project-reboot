package messengers

import (
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

	formattedMessage, err := formatDefaultMessage(message)
	if err != nil {
		return err
	}

	channel, err := discord.session.UserChannelCreate(message.User.AlertDiscordId)
	if err != nil {
		return err
	}
	_, err = discord.session.ChannelMessageSend(channel.ID, formattedMessage)
	if err != nil {
		return err
	}

	return nil
}
