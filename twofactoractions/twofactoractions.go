package twofactoractions

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/jonboulle/clockwork"
)

const CODE_LENGTH = 6

var DEFAULT_CODE_LIFETIME = 2 * time.Minute

var ErrDatabase = errors.New("database error")

func Create(
	actionType string,
	version int,
	expiresAt time.Time,
	data any,
	dbClient *ent.Client,
) (uuid.UUID, string, error) {
	encoded, err := Encode(
		fmt.Sprintf("%v_%v", actionType, version),
		data,
	)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	code := common.CryptoRandomAlphaNum(CODE_LENGTH)
	action, err := dbClient.TwoFactorAction.Create().
		SetType(actionType).
		SetVersion(version).
		SetData(encoded).
		SetExpiresAt(expiresAt).
		SetCode(code).Save(context.Background())
	if err != nil {
		return uuid.UUID{}, code, ErrDatabase
	}

	return action.ID, code, nil
}

var ErrNotFound = errors.New("no action with given ID")
var ErrWrongCode = errors.New("wrong 2FA code")
var ErrExpired = errors.New("action has expired")
var ErrUnknownActionType = errors.New("unknown action type")
var ErrInvalidData = errors.New("invalid action data")

func Confirm(actionID uuid.UUID, code string, db common.DatabaseService, clock clockwork.Clock) error {
	mu := db.TwoFactorActionMutex()
	dbClient := db.Client()
	mu.Lock()

	action, err := dbClient.TwoFactorAction.Get(context.Background(), actionID)
	if err != nil {
		mu.Unlock()
		return ErrNotFound
	}
	if subtle.ConstantTimeCompare([]byte(code), []byte(action.Code)) == 0 {
		mu.Unlock()
		return ErrWrongCode
	}

	err = dbClient.TwoFactorAction.DeleteOne(action).Exec(context.Background())
	mu.Unlock()
	if err != nil {
		return ErrDatabase
	}

	if action.ExpiresAt.Before(clock.Now()) {
		return ErrExpired
	}
	fullType := fmt.Sprintf("%v_%v", action.Type, action.Version)
	actionFunc, ok := actionMap[fullType]
	if !ok {
		return ErrUnknownActionType
	}

	return actionFunc(action)
}
