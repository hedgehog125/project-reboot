package messengers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/common"
)

type discord struct {
	env *common.Env
}

type Discord interface {
	Messenger
}

func NewDiscord(env *common.Env) Discord {
	discord := &discord{
		env: env,
	}

	// Check we can create a session, even though it won't be cached due to the extra complexity
	session, err := discord.getSession()
	if err != nil {
		log.Fatalf("error creating Discord session:\n%v", err)
	}
	err = session.Close()
	if err != nil {
		log.Fatalf("error closing Discord session:\n%v", err)
	}
	return discord
}

func (discord *discord) getSession() (*discordgo.Session, error) {
	return discordgo.New("Bot " + discord.env.DISCORD_TOKEN)
}

func (discord *discord) Id() string {
	return "discord"
}

func (discord *discord) Send(message common.Message) error {
	session, err := discord.getSession()
	if err != nil {
		return err
	}
	err = session.Open()
	if err != nil {
		return err
	}

	// TODO: why does calling close and returning here cause a log?
	// Looks like it's trying to send a heartbeat after closing
	// Maybe need to explicitly stop listening for VC events since they aren't being used?
	defer session.Close()

	formattedMessage, err := formatDefaultMessage(message)
	if err != nil {
		return err
	}

	channel, err := session.UserChannelCreate(message.User.AlertDiscordId)
	if err != nil {
		return err
	}
	_, err = session.ChannelMessageSend(channel.ID, formattedMessage)
	if err != nil {
		return err
	}

	return nil
}
