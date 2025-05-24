package twofactoractions

import (
	"encoding/json"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeEncoding = "encoding"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

var ErrUnknownActionType = common.WrapErrorWithCategory(
	nil, common.ErrTypeTwoFactorAction, "unknown action type",
)

func (registry *Registry) Encode(fullType string, data any) (string, *common.Error) {
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return "", ErrUnknownActionType.AddCategory(ErrTypeEncoding)
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", common.WrapErrorWithCategory(
			err, ErrTypeInvalidData, ErrTypeEncoding,
		)
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType
	err = json.Unmarshal(encoded, &temp)
	if err != nil {
		return "", common.WrapErrorWithCategory(
			err, ErrTypeInvalidData, ErrTypeEncoding,
		)
	}

	return string(encoded), nil
}
