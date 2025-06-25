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
	// TODO: only sending the message once we know we have a format for it prevents the warning, but why does it happen in the first place?
	formattedMessage, commErr := formatDefaultMessage(message)
	if commErr != nil {
		return commErr.AddCategory(ErrTypeSend)
	}

	session, commErr := discord.getSession()
	if commErr != nil {
		return commErr.AddCategory(ErrTypeSend)
	}
	stdErr := session.Open()
	if stdErr != nil {
		return ErrWrapperAPI.Wrap(stdErr).AddCategory(ErrTypeSend)
	}

	// TODO: why does calling close and returning here cause a log?
	// Looks like it's trying to send a heartbeat after closing
	// Maybe need to explicitly stop listening for VC events since they aren't being used?
	defer session.Close()

	channel, stdErr := session.UserChannelCreate(message.User.AlertDiscordId)
	if stdErr != nil {
		return ErrWrapperAPI.Wrap(stdErr).AddCategory(ErrTypeSend)
	}
	_, stdErr = session.ChannelMessageSend(channel.ID, formattedMessage)
	if stdErr != nil {
		return ErrWrapperAPI.Wrap(stdErr).AddCategory(ErrTypeSend)
	}

	return nil
}
