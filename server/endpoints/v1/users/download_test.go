package users_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/common/testcommon"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/server/endpoints/testhelpers"
	"github.com/hedgehog125/project-reboot/server/endpoints/v1/users"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func TestDownload_HappyPath(t *testing.T) {
	t.Parallel()
	// TODO: assert messenger sent message, maybe improve the setup

	clock := clockwork.NewFakeClock()
	app := testhelpers.NewApp(t, &testhelpers.AppOptions{
		Clock: clock,
	})

	username := "alice"
	password := "password123456"
	fileContent := []byte("file content here")
	filename := "data.zip"
	mimeType := "application/zip"

	keySalt := app.Core.GenerateSalt()
	encryptionKey := app.Core.HashPassword(password, keySalt, app.Env.PASSWORD_HASH_SETTINGS)

	encrypted, nonce, wrappedErr := app.Core.Encrypt(fileContent, encryptionKey)
	require.NoError(t, wrappedErr)

	sessionOb, stdErr := dbcommon.WithReadWriteTx(
		t.Context(), app.Database,
		func(tx *ent.Tx, ctx context.Context) (*ent.Session, error) {
			userOb, stdErr := tx.User.Create().
				SetUsername(username).
				SetSessionsValidFrom(clock.Now()).
				SetContent(encrypted).
				SetFileName(filename).
				SetMime(mimeType).
				SetNonce(nonce).
				SetKeySalt(keySalt).
				SetHashTime(app.Env.PASSWORD_HASH_SETTINGS.Time).
				SetHashMemory(app.Env.PASSWORD_HASH_SETTINGS.Memory).
				SetHashThreads(app.Env.PASSWORD_HASH_SETTINGS.Threads).
				Save(ctx)
			if stdErr != nil {
				return nil, stdErr
			}

			authCode := app.Core.RandomAuthCode()
			now := clock.Now()
			validUntil := now.Add(24 * time.Hour)

			sessionOb, stdErr := tx.Session.Create().
				SetUser(userOb).
				SetCreatedAt(clock.Now()).
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
				SetSentAt(clock.Now()).
				SetVersionedMessengerType(app.MockMessenger.VersionedName()).
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

func TestDownload_UndeletedInvalidSession_ReturnsUnauthorizedError(t *testing.T) {
	t.Parallel()

	clock := clockwork.NewFakeClock()
	app := testhelpers.NewApp(t, &testhelpers.AppOptions{
		Clock: clock,
	})

	username := "bob"
	password := "password123456"
	fileContent := []byte("file content here")
	filename := "data.zip"
	mimeType := "application/zip"

	keySalt := core.GenerateSalt()
	encryptionKey := core.HashPassword(password, keySalt, app.Env.PASSWORD_HASH_SETTINGS)

	encrypted, nonce, wrappedErr := core.Encrypt(fileContent, encryptionKey)
	require.NoError(t, wrappedErr)

	sessionOb, stdErr := dbcommon.WithReadWriteTx(
		t.Context(), app.Database,
		func(tx *ent.Tx, ctx context.Context) (*ent.Session, error) {
			now := clock.Now()
			// Set SessionsValidFrom to be in the future
			sessionsValidFrom := now.Add(1 * time.Hour)

			userOb, stdErr := tx.User.Create().
				SetUsername(username).
				SetSessionsValidFrom(clock.Now()).
				SetContent(encrypted).
				SetFileName(filename).
				SetMime(mimeType).
				SetNonce(nonce).
				SetKeySalt(keySalt).
				SetHashTime(app.Env.PASSWORD_HASH_SETTINGS.Time).
				SetHashMemory(app.Env.PASSWORD_HASH_SETTINGS.Memory).
				SetHashThreads(app.Env.PASSWORD_HASH_SETTINGS.Threads).
				SetSessionsValidFrom(sessionsValidFrom).
				Save(ctx)
			if stdErr != nil {
				return nil, stdErr
			}

			authCode := core.RandomAuthCode()
			validUntil := now.Add(24 * time.Hour)

			sessionOb, stdErr := tx.Session.Create().
				SetUser(userOb).
				SetCreatedAt(clock.Now()).
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
				SetSentAt(clock.Now()).
				SetVersionedMessengerType(app.MockMessenger.VersionedName()).
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
	require.Equal(t, http.StatusUnauthorized, respRecorder.Code)
}

func TestDownload_TemporarilyLockedUser_ReturnsUnauthorizedError(t *testing.T) {
	panic("not implemented")
}
func TestDownload_ExpiredTemporarilyLockedUser_AllowsDownload(t *testing.T) {
	panic("not implemented")
}
func TestDownload_PermanentlyLockedUser_ReturnsUnauthorizedError(t *testing.T) {
	panic("not implemented")
}
