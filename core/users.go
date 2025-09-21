package core

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/jonboulle/clockwork"
)

// Doubled because the bytes are represented as base64
const AuthCodeByteLength = 128

func RandomAuthCode() []byte {
	return common.CryptoRandomBytes(AuthCodeByteLength)
}

func SendActiveSessionReminders(ctx context.Context, clock clockwork.Clock, messengers common.MessengerService) *common.Error {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return ErrWrapperSendActiveSessionReminders.Wrap(common.ErrNoTxInContext)
	}

	userObs, stdErr := tx.User.Query().
		WithSessions(func(sessionQuery *ent.SessionQuery) {
			sessionQuery.
				Where(session.ValidUntilGT(clock.Now())).
				Order(ent.Asc(session.FieldValidFrom)).
				Limit(1)
		}).
		All(ctx)
	if stdErr != nil {
		return ErrWrapperSendActiveSessionReminders.Wrap(
			ErrWrapperDatabase.Wrap(stdErr),
		)
	}

	messages := make([]*common.Message, 0, len(userObs))
	for _, userOb := range userObs {
		sessionObs := userOb.Edges.Sessions
		if len(sessionObs) == 0 {
			continue
		}
		sessionOb := sessionObs[0]

		messages = append(messages, &common.Message{
			Type: common.MessageActiveSessionReminder,
			User: userOb,
			Time: sessionOb.ValidFrom,
		})
	}
	commErr := messengers.SendBulk(messages, ctx)
	if commErr != nil {
		return ErrWrapperSendActiveSessionReminders.Wrap(commErr)
	}

	return nil
}

func DeleteExpiredSessions(ctx context.Context, clock clockwork.Clock) *common.Error {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return ErrWrapperDeleteExpiredSessions.Wrap(common.ErrNoTxInContext)
	}

	_, stdErr := tx.Session.Delete().
		Where(session.ValidUntilLTE(clock.Now())).
		Exec(ctx)
	if stdErr != nil {
		return ErrWrapperDeleteExpiredSessions.Wrap(
			ErrWrapperDatabase.Wrap(stdErr),
		)
	}
	return nil
}
