package services

import (
	"context"
	"crypto/subtle"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/twofactoraction"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type TwoFactorActions struct {
	app *common.App
}

func NewTwoFactorActions(app *common.App) *TwoFactorActions {
	return &TwoFactorActions{
		app: app,
	}
}

func (service *TwoFactorActions) Create(
	versionedType string,
	expiresAt time.Time,
	body any,
	ctx context.Context,
) (*ent.TwoFactorAction, string, *common.Error) {
	encoded, commErr := service.app.Jobs.Encode(
		versionedType,
		body,
	)
	if commErr != nil {
		return nil, "", twofactoractions.ErrWrapperCreate.Wrap(commErr)
	}
	jobType, version, commErr := common.ParseVersionedType(versionedType)
	if commErr != nil { // This shouldn't happen because of the Encode call but just in case
		return nil, "", twofactoractions.ErrWrapperCreate.Wrap(commErr)
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return nil, "", twofactoractions.ErrWrapperCreate.Wrap(
			twofactoractions.ErrNoTxInContext,
		)
	}
	code := common.CryptoRandomAlphaNum(twofactoractions.CODE_LENGTH)
	action, err := tx.TwoFactorAction.Create().
		SetType(jobType).
		SetVersion(version).
		SetBody(encoded).
		SetExpiresAt(expiresAt).
		SetCode(code).Save(ctx)
	if err != nil {
		return nil, code, twofactoractions.ErrWrapperCreate.Wrap(
			twofactoractions.ErrWrapperDatabase.Wrap(err),
		)
	}

	return action, code, nil
}
func (service *TwoFactorActions) Confirm(
	actionID uuid.UUID, code string,
	ctx context.Context,
) (*ent.Job, *common.Error) {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(
			twofactoractions.ErrNoTxInContext,
		)
	}
	action, stdErr := tx.TwoFactorAction.Get(ctx, actionID)
	if stdErr != nil {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(
			twofactoractions.ErrNotFound,
		)
	}
	if subtle.ConstantTimeCompare([]byte(code), []byte(action.Code)) == 0 {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(
			twofactoractions.ErrWrongCode,
		)
	}

	stdErr = tx.TwoFactorAction.DeleteOne(action).Exec(ctx)
	if stdErr != nil {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(
			twofactoractions.ErrWrapperDatabase.Wrap(stdErr),
		)
	}

	if action.ExpiresAt.Before(service.app.Clock.Now()) {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(
			twofactoractions.ErrExpired,
		)
	}

	job, commErr := service.app.Jobs.EnqueueEncoded(
		common.GetVersionedType(action.Type, action.Version),
		action.Body,
		ctx,
	)
	if commErr != nil {
		return nil, twofactoractions.ErrWrapperConfirm.Wrap(commErr)
	}
	return job, nil
}
func (service *TwoFactorActions) DeleteExpiredActions(ctx context.Context) *common.Error {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return twofactoractions.ErrWrapperDeleteExpiredActions.Wrap(
			twofactoractions.ErrNoTxInContext,
		)
	}

	_, stdErr := tx.TwoFactorAction.Delete().
		Where(twofactoraction.ExpiresAtLTE(service.app.Clock.Now())).
		Exec(ctx)
	if stdErr != nil {
		return twofactoractions.ErrWrapperDeleteExpiredActions.Wrap(
			twofactoractions.ErrWrapperDatabase.Wrap(stdErr),
		)
	}
	return nil
}
