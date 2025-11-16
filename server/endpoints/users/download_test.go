package users_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/messengers"
	"github.com/hedgehog125/project-reboot/server/endpoints/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/services"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

type mockMessenger struct {
	Name         string
	SentMessages []SentMessage
}
type SentMessage struct {
	Type   string
	UserID int
}

func newMockMessenger(name string) *mockMessenger {
	return &mockMessenger{
		Name:         name,
		SentMessages: []SentMessage{},
	}
}

func (mockMessenger *mockMessenger) Register(registry *messengers.Registry) {
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
func (mockMessenger *mockMessenger) VersionedName() string {
	return common.GetVersionedType(mockMessenger.Name, 1)
}

func TestDownload_HappyPath(t *testing.T) {
	t.Parallel()
	// TODO: assert messenger sent message, maybe improve the setup

	clock := clockwork.NewFakeClock()
	env := testcommon.DefaultEnv()
	/*
		If you want to test performance:
			env.PASSWORD_HASH_SETTINGS = &common.PasswordHashSettings{
				Time:    5,
				Memory:  1 * 1024 * 1024, // 1 GiB
				Threads: 4,
			}
	*/
	app := &common.App{
		Env:   env,
		Clock: clock,
	}
	{
		logger := services.NewLogger(app)
		app.Logger = logger
		slog.SetDefault(logger.Logger)
	}
	app.RateLimiter = services.NewRateLimiter(app)
	app.Core = services.NewCore(app)
	db := testcommon.CreateDB()
	defer db.Shutdown()
	app.Database = db
	app.KeyValue = services.NewKeyValue(app)
	app.KeyValue.Init()
	mockMessenger := newMockMessenger("MOCK_MESSENGER_1")
	{
		messengerService := services.NewMessengers(app, mockMessenger.Register)
		app.Messengers = messengerService
		app.Jobs = services.NewJobs(app, messengerService.RegisterJobs)
	}
	app.Server = services.NewServer(app)

	app.Jobs.Start()
	defer app.Jobs.Shutdown()

	username := "alice"
	password := "password123456"
	fileContent := []byte("file content here")
	filename := "data.zip"
	mimeType := "application/zip"

	keySalt := core.GenerateSalt()
	encryptionKey := core.HashPassword(password, keySalt, env.PASSWORD_HASH_SETTINGS)

	encrypted, nonce, wrappedErr := core.Encrypt(fileContent, encryptionKey)
	require.NoError(t, wrappedErr)

	sessionOb, stdErr := dbcommon.WithReadWriteTx(
		t.Context(), db,
		func(tx *ent.Tx, ctx context.Context) (*ent.Session, error) {
			userOb, stdErr := tx.User.Create().
				SetUsername(username).
				SetContent(encrypted).
				SetFileName(filename).
				SetMime(mimeType).
				SetNonce(nonce).
				SetKeySalt(keySalt).
				SetHashTime(env.PASSWORD_HASH_SETTINGS.Time).
				SetHashMemory(env.PASSWORD_HASH_SETTINGS.Memory).
				SetHashThreads(env.PASSWORD_HASH_SETTINGS.Threads).
				Save(ctx)
			if stdErr != nil {
				return nil, stdErr
			}

			authCode := core.RandomAuthCode()
			now := clock.Now()
			validUntil := now.Add(24 * time.Hour)

			sessionOb, stdErr := tx.Session.Create().
				SetUser(userOb).
				SetCode(authCode).
				SetValidFrom(now).
				SetValidUntil(validUntil).
				SetUserAgent("test-agent").
				SetIP("127.0.0.1").
				Save(ctx)
			if stdErr != nil {
				return sessionOb, stdErr
			}

			_, stdErr = tx.LoginAlert.Create().
				SetSession(sessionOb).
				SetTime(clock.Now()).
				SetVersionedMessengerType(mockMessenger.VersionedName()).
				SetConfirmed(true).
				Save(ctx)
			return sessionOb, stdErr
		},
	)
	require.NoError(t, stdErr)

	respRecorder := testcommon.Post(
		t, app.Server,
		"/api/v1/users/download",
		users.DownloadPayload{
			Username:          username,
			Password:          password,
			AuthorizationCode: base64.StdEncoding.EncodeToString(sessionOb.Code),
		},
	)
	testcommon.AssertJSONResponse(
		t, respRecorder,
		http.StatusOK,
		&users.DownloadResponse{
			Errors:                      []servercommon.ErrorDetail{},
			AuthorizationCodeValidFrom:  nil,
			AuthorizationCodeValidUntil: nil,
			Content:                     fileContent,
			Filename:                    filename,
			Mime:                        mimeType,
		},
	)
	require.Equal(t, http.StatusOK, respRecorder.Code)
}
