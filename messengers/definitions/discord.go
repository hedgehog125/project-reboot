package definitions

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers"
)

type Discord1Body struct {
	UserID           string `json:"userID"`
	FormattedMessage string `json:"formattedMessage"`
}

var ErrWrapperDiscord = common.NewDynamicErrorWrapper(func(err error) *common.Error {
	commErr := common.ErrWrapperAPI.Wrap(err)
	if commErr == nil {
		return nil
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return commErr.
			ConfigureRetries(10, time.Second*5, 1.5).
			AddDebugValue(common.DebugValue{
				Name: "retried url.Error",
			})
	}
	var rateLimitErr *discordgo.RateLimitError
	if errors.As(err, &rateLimitErr) {
		return commErr.
			ConfigureRetries(3, max(rateLimitErr.RetryAfter, 5*time.Second), 1).
			AddDebugValue(common.DebugValue{
				Name:    "retried discordgo.RateLimitError",
				Message: fmt.Sprintf("RetryAfter: %v", rateLimitErr.RetryAfter),
			})
	}

	return commErr
})

func Discord1(app *common.App) *messengers.Definition {
	getSession := func() (*discordgo.Session, *common.Error) {
		session, err := discordgo.New("Bot " + app.Env.DISCORD_TOKEN)
		if err != nil {
			return nil, ErrWrapperDiscord.Wrap(err)
		}

		session.ShouldRetryOnRateLimit = false
		return session, nil
	}
	session, commErr := getSession()
	if commErr != nil {
		log.Fatalf("error creating startup test Discord session:\n%v", commErr)
	}
	stdErr := session.Close()
	if stdErr != nil {
		log.Fatalf("error closing startup test Discord session:\n%v", stdErr)
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
			commErr := jobCtx.Decode(body)
			if commErr != nil {
				return commErr
			}

			session, commErr := getSession()
			if commErr != nil {
				return commErr
			}
			stdErr := session.Open()
			if stdErr != nil {
				return ErrWrapperDiscord.Wrap(stdErr)
			}

			// TODO: why does calling close and returning here cause a log?
			// Looks like it's trying to send a heartbeat after closing
			// Maybe need to explicitly stop listening for VC events since they aren't being used?
			defer func() {
				stdErr := session.Close()
				if stdErr != nil {
					jobCtx.Logger.Warn("error closing Discord session", "error", stdErr)
				}
			}()

			channel, stdErr := session.UserChannelCreate(body.UserID)
			if stdErr != nil {
				return ErrWrapperDiscord.Wrap(stdErr)
			}
			_, stdErr = session.ChannelMessageSend(channel.ID, body.FormattedMessage)
			if stdErr != nil {
				return ErrWrapperDiscord.Wrap(stdErr)
			}

			return nil
		},
	}
}
