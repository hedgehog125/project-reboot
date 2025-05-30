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
	session, sessionErr := discord.getSession()
	if sessionErr != nil {
		log.Fatalf("error creating Discord session:\n%v", sessionErr)
	}
	closeErr := session.Close()
	if closeErr != nil {
		log.Fatalf("error closing Discord session:\n%v", closeErr)
	}
	return discord
}

func (discord *discord) getSession() (*discordgo.Session, *common.Error) {
	session, err := discordgo.New("Bot " + discord.env.DISCORD_TOKEN)
	if err != nil {
		return nil, ErrWrapperAPI.Wrap(err)
	}
	return session, nil
}

func (discord *discord) Id() string {
	return "discord"
}

func (discord *discord) Send(message common.Message) *common.Error {
	session, sessionErr := discord.getSession()
	if sessionErr != nil {
		return sessionErr.AddCategory(ErrTypeSend)
	}
	openErr := session.Open()
	if openErr != nil {
		return ErrWrapperAPI.Wrap(openErr).AddCategory(ErrTypeSend)
	}

	// TODO: why does calling close and returning here cause a log?
	// Looks like it's trying to send a heartbeat after closing
	// Maybe need to explicitly stop listening for VC events since they aren't being used?
	defer session.Close()

	formattedMessage, formatErr := formatDefaultMessage(message)
	if formatErr != nil {
		return formatErr.AddCategory(ErrTypeSend)
	}

	channel, createErr := session.UserChannelCreate(message.User.AlertDiscordId)
	if createErr != nil {
		return ErrWrapperAPI.Wrap(createErr).AddCategory(ErrTypeSend)
	}
	_, sendErr := session.ChannelMessageSend(channel.ID, formattedMessage)
	if sendErr != nil {
		return ErrWrapperAPI.Wrap(sendErr).AddCategory(ErrTypeSend)
	}

	return nil
}
