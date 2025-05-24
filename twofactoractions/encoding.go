package twofactoractions

import (
	"encoding/json"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeInvalidData = "invalid data"
)

// TODO: adapt other error constants to this pattern <=============
// TODO: does this work with errors.Is?
var ErrUnknownActionType = common.WrapErrorWithCategory(
	nil, "unknown action type",
)

func (registry *Registry) Encode(fullType string, data any) (string, error) {
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return "", ErrUnknownActionType
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", common.WrapErrorWithCategory(
			err, ErrTypeInvalidData,
		)
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType
	err = json.Unmarshal(encoded, &temp)
	if err != nil {
		return "", common.WrapErrorWithCategory(
			err, ErrTypeInvalidData,
		)
	}

	return string(encoded), nil
}
