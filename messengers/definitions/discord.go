package definitions

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers"
)

type Discord1Body struct {
	UserID           string `json:"userID"`
	FormattedMessage string `json:"formattedMessage"`
}

func Discord1(app *common.App) *messengers.Definition {
	getSession := func() (*discordgo.Session, *common.Error) {
		session, err := discordgo.New("Bot " + app.Env.DISCORD_TOKEN)
		if err != nil {
			return nil, common.ErrWrapperAPI.Wrap(err)
		}
		return session, nil
	}
	session, commErr := getSession()
	if commErr != nil {
		log.Fatalf("error creating Discord session:\n%v", commErr)
	}
	stdErr := session.Close()
	if stdErr != nil {
		log.Fatalf("error closing Discord session:\n%v", stdErr)
	}

	return &messengers.Definition{
		ID:      "discord",
		Version: 1,
		Prepare: func(message *common.Message) (any, error) {
			if message.User.AlertDiscordId == "" {
				return nil, messengers.ErrNoContactForUser.Clone()
			}
			formattedMessage, commErr := messengers.FormatDefaultMessage(message)
			if commErr != nil {
				return nil, commErr
			}

			return &Discord1Body{
				UserID:           message.User.AlertDiscordId,
				FormattedMessage: formattedMessage,
			}, nil
		},
		BodyType: &Discord1Body{},
		Handler: func(jobCtx *jobs.Context) error {
			body := &Discord1Body{}
			jobErr := jobCtx.Decode(body)
			if jobErr != nil {
				return jobErr
			}

			session, commErr := getSession()
			if commErr != nil {
				return commErr
			}
			stdErr := session.Open()
			if stdErr != nil {
				return common.ErrWrapperAPI.Wrap(stdErr)
			}

			// TODO: why does calling close and returning here cause a log?
			// Looks like it's trying to send a heartbeat after closing
			// Maybe need to explicitly stop listening for VC events since they aren't being used?
			defer func() {
				stdErr := session.Close()
				if stdErr != nil {
					fmt.Printf("warning: error closing Discord session:\n%v\n", stdErr)
				}
			}()

			channel, stdErr := session.UserChannelCreate(body.UserID)
			if stdErr != nil {
				return common.ErrWrapperAPI.Wrap(stdErr)
			}
			_, stdErr = session.ChannelMessageSend(channel.ID, body.FormattedMessage)
			if stdErr != nil {
				return common.ErrWrapperAPI.Wrap(stdErr)
			}

			return nil
		},
	}
}
