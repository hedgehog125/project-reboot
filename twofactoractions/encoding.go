package twofactoractions

import (
	"encoding/json"

	"github.com/hedgehog125/project-reboot/common"
)

func (registry *Registry) Encode(fullType string, data any) (string, *common.Error) {
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return "", ErrUnknownActionType.AddCategory(ErrTypeEncode)
	}

	encoded, stdErr := json.Marshal(data)
	if stdErr != nil {
		return "", ErrWrapperInvalidData.Wrap(stdErr).AddCategory(ErrTypeEncode)
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType()
	stdErr = json.Unmarshal(encoded, temp)
	if stdErr != nil {
		return "", ErrWrapperInvalidData.Wrap(stdErr).AddCategory(ErrTypeEncode)
	}

	return string(encoded), nil
}
