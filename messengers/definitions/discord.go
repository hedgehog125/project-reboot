package definitions

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/messengers"
)

type Discord1Body struct {
	UserID           string `json:"userID"`
	FormattedMessage string `json:"formattedMessage"`
}

var ErrWrapperDiscord = common.NewDynamicErrorWrapper(func(err error) common.WrappedError {
	wrappedErr := common.ErrWrapperAPI.Wrap(err)
	if wrappedErr == nil {
		return nil
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		wrappedErr.ConfigureRetriesMut(10, time.Second*5, 1.5)
		wrappedErr.AddDebugValuesMut(common.DebugValue{
			Name: "retried url.Error",
		})
		return wrappedErr
	}
	var rateLimitErr *discordgo.RateLimitError
	if errors.As(err, &rateLimitErr) {
		wrappedErr.ConfigureRetriesMut(3, max(rateLimitErr.RetryAfter, 5*time.Second), 1)
		wrappedErr.AddDebugValuesMut(common.DebugValue{
			Name:    "retried discordgo.RateLimitError",
			Message: fmt.Sprintf("RetryAfter: %v", rateLimitErr.RetryAfter),
		})
		return wrappedErr
	}

	return wrappedErr
})

func Discord1(app *common.App) *messengers.Definition {
	getSession := func() (*discordgo.Session, common.WrappedError) {
		session, err := discordgo.New("Bot " + app.Env.DISCORD_TOKEN)
		if err != nil {
			return nil, ErrWrapperDiscord.Wrap(err)
		}

		session.ShouldRetryOnRateLimit = false
		return session, nil
	}
	session, wrappedErr := getSession()
	if wrappedErr != nil {
		log.Fatalf("error creating startup test Discord session:\n%v", wrappedErr)
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
			formattedMessage, wrappedErr := messengers.FormatDefaultMessage(message)
			if wrappedErr != nil {
				return nil, wrappedErr
			}

			return &Discord1Body{
				UserID:           message.User.AlertDiscordId,
				FormattedMessage: formattedMessage,
			}, nil
		},
		BodyType: &Discord1Body{},
		Handler: func(messengerCtx *messengers.Context) error {
			body := &Discord1Body{}
			wrappedErr := messengerCtx.Decode(body)
			if wrappedErr != nil {
				return wrappedErr
			}

			session, wrappedErr := getSession()
			if wrappedErr != nil {
				return wrappedErr
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
					messengerCtx.Logger.Warn("error closing Discord session", "error", stdErr)
				}
			}()

			// TODO: it's probably worth caching this to reduce how often we're rate limited
			channel, stdErr := session.UserChannelCreate(body.UserID)
			if stdErr != nil {
				return ErrWrapperDiscord.Wrap(stdErr)
			}
			_, stdErr = session.ChannelMessageSend(channel.ID, body.FormattedMessage)
			if stdErr != nil {
				return ErrWrapperDiscord.Wrap(stdErr)
			}

			messengerCtx.ConfirmSent()
			return nil
		},
	}
}
