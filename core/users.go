package core

import (
	"context"
	"fmt"
	"math"
	"slices"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/ent/session"
	"github.com/NicoClack/cryptic-stash/ent/user"
	"github.com/jonboulle/clockwork"
)

// Doubled because the bytes are represented as base64
const AuthCodeByteLength = 128

func RandomAuthCode() []byte {
	return common.CryptoRandomBytes(AuthCodeByteLength)
}

func SendActiveSessionReminders(
	ctx context.Context,
	clock clockwork.Clock,
	messengers common.MessengerService,
) common.WrappedError {
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
		sessionIDs := make([]int, 0, len(sessionObs))
		for _, sessionOb := range sessionObs {
			sessionIDs = append(sessionIDs, sessionOb.ID)
		}

		messages = append(messages, &common.Message{
			Type:       common.MessageActiveSessionReminder,
			User:       userOb,
			Time:       sessionOb.ValidFrom,
			SessionIDs: sessionIDs,
		})
	}
	wrappedErr := messengers.SendBulk(messages, ctx)
	if wrappedErr != nil {
		return ErrWrapperSendActiveSessionReminders.Wrap(wrappedErr)
	}

	return nil
}

func DeleteExpiredSessions(ctx context.Context, clock clockwork.Clock) common.WrappedError {
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

func InvalidateUserSessions(userID int, ctx context.Context, clock clockwork.Clock) common.WrappedError {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return ErrWrapperInvalidateUserSessions.Wrap(common.ErrNoTxInContext)
	}

	_, stdErr := tx.Session.Delete().
		Where(session.HasUserWith(user.ID(userID))).
		Exec(ctx)
	if stdErr != nil {
		return ErrWrapperInvalidateUserSessions.Wrap(
			ErrWrapperDatabase.Wrap(stdErr),
		)
	}
	stdErr = tx.User.UpdateOneID(userID).
		SetSessionsValidFrom(clock.Now()).
		Exec(ctx)
	if stdErr != nil {
		return ErrWrapperInvalidateUserSessions.Wrap(
			ErrWrapperDatabase.Wrap(stdErr),
		)
	}
	return nil
}

func IsUserSufficientlyNotified(
	sessionOb *ent.Session,
	messengers common.MessengerService,
	logger common.Logger,
	clock clockwork.Clock, env *common.Env,
) bool {
	logger = logger.With(
		"sessionID", sessionOb.ID,
		"userID", sessionOb.Edges.User.ID,
	)

	allLoginAlerts := slices.Clone(sessionOb.Edges.LoginAlerts)
	groupedLoginAlerts := make(map[string][]*ent.LoginAlert)
	for _, loginAlert := range allLoginAlerts {
		groupedLoginAlerts[loginAlert.VersionedMessengerType] = append(
			groupedLoginAlerts[loginAlert.VersionedMessengerType],
			loginAlert,
		)
	}
	messengerTypes := messengers.GetConfiguredMessengerTypes(sessionOb.Edges.User)
	earliestValidTime := clock.Now().Add(-env.ACTIVE_SESSION_REMINDER_INTERVAL)
	successfulMessengerTypes := []string{}
	// Ignore the supplemental messengers when assessing this
	coreMessengerTypeCount := 0
	for _, messengerType := range messengerTypes {
		messengerDef, ok := messengers.GetPublicDefinition(messengerType)
		if !ok {
			panic(fmt.Sprintf("IsUserSufficientlyNotified: no messenger definition for %s", messengerType))
		}
		if messengerDef.IsSupplemental {
			continue
		}
		coreMessengerTypeCount++

		loginAlerts := groupedLoginAlerts[messengerType]
		confirmedLoginAlerts := []*ent.LoginAlert{}
		for _, alert := range loginAlerts {
			if alert.Confirmed {
				confirmedLoginAlerts = append(confirmedLoginAlerts, alert)
			}
		}

		if len(confirmedLoginAlerts) < env.MIN_SUCCESSFUL_MESSAGE_COUNT {
			logger.Warn(
				"user was not sufficiently notified by one of their configured messengers because it "+
					"didn't successfully send and confirm enough login alerts",
				"messengerType",
				messengerType,
				"loginAlertCount",
				len(loginAlerts),
				"confirmedLoginAlertCount",
				len(confirmedLoginAlerts),
			)
			continue
		}
		mostRecentConfirmedAlert := &ent.LoginAlert{}
		for _, alert := range confirmedLoginAlerts {
			if alert.SentAt.After(mostRecentConfirmedAlert.SentAt) {
				mostRecentConfirmedAlert = alert
			}
		}
		if mostRecentConfirmedAlert.SentAt.Before(earliestValidTime) {
			logger.Warn(
				"user was not sufficiently notified by one of their configured messengers because "+
					"its most recent confirmed alert was too old. are jobs still running? are some messengers failing?",
				"messengerType",
				messengerType,
				"mostRecentAlertTime",
				mostRecentConfirmedAlert.SentAt,
				"earliestValidTime",
				earliestValidTime,
			)
			continue
		}

		successfulMessengerTypes = append(successfulMessengerTypes, messengerType)
	}

	minSuccessfulMessengers := max(int(
		math.Ceil(float64(coreMessengerTypeCount)/float64(2)),
	), 1)
	if len(successfulMessengerTypes) < minSuccessfulMessengers {
		logger.Warn(
			"user was not sufficiently notified because not enough of their core configured messengers "+
				"successfully sent login alerts",
			"configuredMessengerTypes",
			messengerTypes,
			"successfulMessengerTypes",
			successfulMessengerTypes,
			"minSuccessfulMessengers",
			minSuccessfulMessengers,
		)
		return false
	}

	logger.Info(
		"user was sufficiently notified",
		"configuredMessengerTypes", messengerTypes,
		"successfulMessengerTypes", successfulMessengerTypes,
		"minSuccessfulMessengers", minSuccessfulMessengers,
	)
	return true
}

func IsUserLocked(userOb *ent.User, clock clockwork.Clock) bool {
	if userOb.Locked {
		return true
	}
	if userOb.LockedUntil == nil {
		return false
	}
	return clock.Now().Before(*userOb.LockedUntil)
}
