package twofactoractions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
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

var ErrUnknownActionType = errors.New("unknown action type")
var ErrInvalidData = errors.New("invalid action data")

func Confirm(action *ent.TwoFactorAction) error {
	// TODO: move expiry checking to this
	// TODO: delete as soon as it's read to prevent running twice?

	fullType := fmt.Sprintf("%v_%v", action.Type, action.Version)
	actionFunc, ok := actionMap[fullType]
	if !ok {
		return ErrUnknownActionType
	}

	return actionFunc(action)
}
