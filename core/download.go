package core

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/jonboulle/clockwork"
)

func Download(
	username string,
	password string,
	givenAuthCode []byte,
	encryptor common.EncryptionService,
	db common.DatabaseService,
	logger common.Logger,
	clock clockwork.Clock,
	env *common.Env,
	ctx context.Context,
) ([]byte, *ent.Session, common.WrappedError) {
	sessionOb, stdErr := dbcommon.WithReadTx(
		ctx, db,
		func(tx *ent.Tx, ctx context.Context) (*ent.Session, error) {
			sessionOb, stdErr := tx.Session.Query().
				Where(session.And(session.HasUserWith(user.Username(username)), session.Code(givenAuthCode))).
				WithUser().
				First(ctx)
			if stdErr != nil {
				return nil, servercommon.SendUnauthorizedIfNotFound(stdErr)
			}
			if clock.Now().After(sessionOb.ValidUntil) {
				stdErr := tx.Session.DeleteOneID(sessionOb.ID).Exec(ctx)
				if stdErr != nil {
					logger.Error(
						"unable to delete expired session",
						"error",
						stdErr,
						"sessionID",
						sessionOb.ID,
					)
				}
				return nil, servercommon.NewUnauthorizedError()
			}
			return sessionOb, nil
		},
	)
	if stdErr != nil {
		// TODO: wrap
		return stdErr
	}
	if clock.Now().Before(sessionOb.ValidFrom) {
		return nil, sessionOb, ErrWrapperDownload.Wrap(
			ErrAuthorizationCodeNotValidYet,
		)
	}

	userOb := sessionOb.Edges.User
	encryptionKey := encryptor.HashPassword(
		password,
		userOb.KeySalt,
		&common.PasswordHashSettings{
			Time:    userOb.HashTime,
			Memory:  userOb.HashMemory,
			Threads: userOb.HashThreads,
		},
	)
	decrypted, wrappedErr := encryptor.Decrypt(userOb.Content, encryptionKey, userOb.Nonce)
	if wrappedErr != nil {
		return nil, sessionOb, ErrWrapperDownload.Wrap(
			wrappedErr,
		)
	}

	stdErr = dbcommon.WithWriteTx(
		ctx, db,
		func(tx *ent.Tx, ctx context.Context) error {
			stdErr := tx.Session.UpdateOneID(sessionOb.ID).
				SetValidUntil(clock.Now().Add(env.USED_AUTH_CODE_VALID_FOR)).
				Exec(ctx)
			if stdErr != nil {
				return stdErr
			}
			_, _, wrappedErr := app.Messengers.SendUsingAll(
				&common.Message{
					Type: common.MessageDownload,
					User: userOb,
				},
				ctx,
			)
			return wrappedErr
		},
	)
	if stdErr != nil {
		return stdErr
	}

	return decrypted, sessionOb, nil
}
