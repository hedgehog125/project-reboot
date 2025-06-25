package twofactoractions

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
)

const CODE_LENGTH = 9

var DEFAULT_CODE_LIFETIME = 2 * time.Minute
var MAX_ACTION_RUN_TIME = 15 * time.Second

// TODO: move functions to separate files

func (registry *Registry) Create(
	actionType string,
	version int,
	expiresAt time.Time,
	data any,
) (uuid.UUID, string, *common.Error) {
	encoded, encodeErr := registry.Encode(
		GetVersionedType(actionType, version),
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

	action, stdErr := dbClient.TwoFactorAction.Get(context.Background(), actionID)
	if stdErr != nil {
		mu.Unlock()
		return ErrNotFound.AddCategory(ErrTypeConfirm)
	}
	if subtle.ConstantTimeCompare([]byte(code), []byte(action.Code)) == 0 {
		mu.Unlock()
		return ErrWrongCode.AddCategory(ErrTypeConfirm)
	}

	stdErr = dbClient.TwoFactorAction.DeleteOne(action).Exec(context.Background())
	mu.Unlock()
	if stdErr != nil {
		return ErrWrapperDatabase.Wrap(stdErr).AddCategory(ErrTypeConfirm)
	}

	if action.ExpiresAt.Before(registry.App.Clock.Now()) {
		return ErrExpired.AddCategory(ErrTypeConfirm)
	}
	fullType := GetVersionedType(action.Type, action.Version)
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return ErrUnknownActionType.AddCategory(ErrTypeConfirm)
	}

	parsed := actionDef.BodyType
	stdErr = json.Unmarshal([]byte(action.Data), &parsed)
	// TODO: data is an interface {} (map[string]any)
	// Maybe just do the JSON decoding in the action handler?
	if stdErr != nil {
		return ErrWrapperInvalidData.Wrap(stdErr).AddCategory(ErrTypeConfirm)
	}

	// TODO: standardise this error
	ctx, cancel := context.WithTimeoutCause(context.Background(), MAX_ACTION_RUN_TIME, errors.New("action run time exceeded"))
	defer cancel()
	// TODO: how should this be coordinated with the service?
	// TODO: should also stop new actions from being run during shutdown, that way the server service can be shut down at the same time. This service will just need a slightly longer timeout than it

	return actionDef.Handler(&Action[any]{
		Definition: &actionDef,
		ExpiresAt:  action.ExpiresAt,
		Context:    ctx,
		Body:       &parsed,
	})
}

func GetVersionedType(actionType string, version int) string {
	return fmt.Sprintf("%v_%v", actionType, version)
}
