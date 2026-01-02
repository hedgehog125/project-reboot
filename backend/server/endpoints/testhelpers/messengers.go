package testhelpers

import (
	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/messengers"
)

type MockMessenger struct {
	Name         string
	SentMessages []SentMessage
}
type SentMessage struct {
	Type   string
	UserID int
}

func NewMockMessenger(name string) *MockMessenger {
	return &MockMessenger{
		Name:         name,
		SentMessages: []SentMessage{},
	}
}

func (mockMessenger *MockMessenger) Register(registry *messengers.Registry) {
	type Body struct {
		Message SentMessage
	}
	registry.Register(&messengers.Definition{
		ID:      mockMessenger.Name,
		Version: 1,
		Prepare: func(message *common.Message) (any, error) {
			return &Body{
				Message: SentMessage{
					Type:   string(message.Type),
					UserID: message.User.ID,
				},
			}, nil
		},
		BodyType: &Body{},
		Handler: func(messengerCtx *messengers.Context) error {
			body := Body{}
			wrappedErr := messengerCtx.Decode(&body)
			if wrappedErr != nil {
				return wrappedErr
			}

			mockMessenger.SentMessages = append(mockMessenger.SentMessages, body.Message)
			messengerCtx.ConfirmSent()
			return nil
		},
	})
}
func (mockMessenger *MockMessenger) VersionedName() string {
	return common.GetVersionedType(mockMessenger.Name, 1)
}
