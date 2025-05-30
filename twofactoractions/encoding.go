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

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", ErrWrapperInvalidData.Wrap(err).AddCategory(ErrTypeEncode)
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType
	err = json.Unmarshal(encoded, &temp)
	if err != nil {
		return "", ErrWrapperInvalidData.Wrap(err).AddCategory(ErrTypeEncode)
	}

	return string(encoded), nil
}
