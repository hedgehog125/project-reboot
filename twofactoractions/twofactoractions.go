package twofactoractions

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
)

const CODE_LENGTH = 9

var DEFAULT_CODE_LIFETIME = 2 * time.Minute

// TODO: move functions to separate files

func (registry *Registry) Create(
	actionType string,
	version int,
	expiresAt time.Time,
	data any,
) (uuid.UUID, string, *common.Error) {
	encoded, encodeErr := registry.Encode(
		GetFullType(actionType, version),
		data,
	)
	if encodeErr != nil {
		return uuid.UUID{}, "", encodeErr.AddCategory(ErrTypeCreate)
	}

	code := common.CryptoRandomAlphaNum(CODE_LENGTH)
	dbClient := registry.App.Database.Client()
	action, err := dbClient.TwoFactorAction.Create().
		SetType(actionType).
		SetVersion(version).
		SetData(encoded).
		SetExpiresAt(expiresAt).
		SetCode(code).Save(context.Background())
	if err != nil {
		return uuid.UUID{}, code, ErrWrapperDatabase.Wrap(err).AddCategory(ErrTypeCreate)
	}

	return action.ID, code, nil
}

func (registry *Registry) Confirm(actionID uuid.UUID, code string) *common.Error {
	mu := registry.App.Database.TwoFactorActionMutex()
	dbClient := registry.App.Database.Client()
	mu.Lock()

	action, err := dbClient.TwoFactorAction.Get(context.Background(), actionID)
	if err != nil {
		mu.Unlock()
		return ErrNotFound.AddCategory(ErrTypeConfirm)
	}
	if subtle.ConstantTimeCompare([]byte(code), []byte(action.Code)) == 0 {
		mu.Unlock()
		return ErrWrongCode.AddCategory(ErrTypeConfirm)
	}

	err = dbClient.TwoFactorAction.DeleteOne(action).Exec(context.Background())
	mu.Unlock()
	if err != nil {
		return ErrWrapperDatabase.Wrap(err).AddCategory(ErrTypeConfirm)
	}

	if action.ExpiresAt.Before(registry.App.Clock.Now()) {
		return ErrExpired.AddCategory(ErrTypeConfirm)
	}
	fullType := GetFullType(action.Type, action.Version)
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return ErrUnknownActionType.AddCategory(ErrTypeConfirm)
	}

	parsed := actionDef.BodyType
	err = json.Unmarshal([]byte(action.Data), &parsed)
	if err != nil {
		return ErrWrapperInvalidData.Wrap(err).AddCategory(ErrTypeConfirm)
	}

	return actionDef.Handler(&Action[any]{
		Definition: &actionDef,
		ExpiresAt:  action.ExpiresAt,
		Body:       &parsed,
	})
}

func GetFullType(actionType string, version int) string {
	return fmt.Sprintf("%v_%v", actionType, version)
}
