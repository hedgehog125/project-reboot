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

const (
	ErrTypeConfirm = "confirm"
)

// TODO: move functions to separate files

func (registry *Registry) Create(
	actionType string,
	version int,
	expiresAt time.Time,
	data any,
) (uuid.UUID, string, error) {
	encoded, encodeErr := registry.Encode(
		GetFullType(actionType, version),
		data,
	)
	if encodeErr != nil {
		return uuid.UUID{}, "", encodeErr
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
		return uuid.UUID{}, code, common.WrapErrorWithCategories(err, common.ErrTypeDatabase, common.ErrTypeTwoFactorAction)
	}

	return action.ID, code, nil
}

var ErrNotFound = common.NewErrorWithCategories(
	"no action with given ID", common.ErrTypeTwoFactorAction,
)
var ErrWrongCode = common.NewErrorWithCategories(
	"wrong 2FA code", common.ErrTypeTwoFactorAction,
)
var ErrExpired = common.NewErrorWithCategories(
	"action has expired", common.ErrTypeTwoFactorAction,
)

func (registry *Registry) Confirm(actionID uuid.UUID, code string) error {
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
		return common.WrapErrorWithCategories(
			err, common.ErrTypeDatabase,
			ErrTypeConfirm, common.ErrTypeTwoFactorAction,
		)
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
		return common.WrapErrorWithCategories(
			err, ErrTypeInvalidData, ErrTypeConfirm,
			common.ErrTypeTwoFactorAction,
		)
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
