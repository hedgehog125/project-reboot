package services

import (
	"context"
	"crypto/subtle"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/jobs/jobscommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

func NewTwoFactorAction(app *common.App) common.TwoFactorActionService {
	return &twoFactorActionService{
		app: app,
	}
}

type twoFactorActionService struct {
	app *common.App
}

// TODO: define these errors and constants in twofactoractions package
func (service *twoFactorActionService) Create(
	versionedType string,
	expiresAt time.Time,
	data any,
	ctx context.Context,
) (uuid.UUID, string, *common.Error) {
	encoded, commErr := service.app.Jobs.Encode(
		versionedType,
		data,
	)
	if commErr != nil {
		return uuid.UUID{}, "", twofactoractions.ErrWrapperCreate.Wrap(commErr)
	}
	jobType, version, commErr := jobscommon.ParseVersionedType(versionedType)
	if commErr != nil { // This shouldn't happen because of the Encode call but just in case
		return uuid.UUID{}, "", twofactoractions.ErrWrapperCreate.Wrap(commErr)
	}

	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return uuid.UUID{}, "", twofactoractions.ErrNoTxInContext.AddCategory(twofactoractions.ErrTypeCreate)
	}
	code := common.CryptoRandomAlphaNum(twofactoractions.CODE_LENGTH)
	action, err := tx.TwoFactorAction.Create().
		SetType(jobType).
		SetVersion(version).
		SetData(encoded).
		SetExpiresAt(expiresAt).
		SetCode(code).Save(ctx)
	if err != nil {
		return uuid.UUID{}, code, twofactoractions.ErrWrapperDatabase.Wrap(err).AddCategory(twofactoractions.ErrTypeCreate)
	}

	return action.ID, code, nil
}
func (service *twoFactorActionService) Confirm(
	actionID uuid.UUID, code string,
	ctx context.Context,
) (uuid.UUID, *common.Error) {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return uuid.UUID{}, twofactoractions.ErrNoTxInContext.AddCategory(twofactoractions.ErrTypeConfirm)
	}
	action, stdErr := tx.TwoFactorAction.Get(ctx, actionID)
	if stdErr != nil {
		return uuid.UUID{}, twofactoractions.ErrNotFound.AddCategory(twofactoractions.ErrTypeConfirm)
	}
	if subtle.ConstantTimeCompare([]byte(code), []byte(action.Code)) == 0 {
		return uuid.UUID{}, twofactoractions.ErrWrongCode.AddCategory(twofactoractions.ErrTypeConfirm)
	}

	stdErr = tx.TwoFactorAction.DeleteOne(action).Exec(ctx)
	if stdErr != nil {
		return uuid.UUID{}, twofactoractions.ErrWrapperDatabase.Wrap(stdErr).AddCategory(twofactoractions.ErrTypeConfirm)
	}

	if action.ExpiresAt.Before(service.app.Clock.Now()) {
		return uuid.UUID{}, twofactoractions.ErrExpired.AddCategory(twofactoractions.ErrTypeConfirm)
	}

	jobID, commErr := service.app.Jobs.Enqueue(
		jobscommon.GetVersionedType(action.Type, action.Version),
		action.Data,
		ctx,
	)
	if commErr != nil {
		return uuid.UUID{}, twofactoractions.ErrWrapperConfirm.Wrap(commErr)
	}
	return jobID, nil
}
