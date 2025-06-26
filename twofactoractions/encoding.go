package twofactoractions

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hedgehog125/project-reboot/common"
)

func (registry *Registry) Encode(versionedType string, data any) (string, *common.Error) {
	actionDef, ok := registry.actions[versionedType]
	if !ok {
		return "", ErrUnknownActionType.AddCategory(ErrTypeEncode)
	}

	dataType := reflect.TypeOf(data)
	if dataType != actionDef.reflectedBodyType {
		return "", ErrWrapperInvalidData.Wrap(
			fmt.Errorf("data type %s isn't the expected type %s", dataType, actionDef.reflectedBodyType),
		).AddCategory(ErrTypeEncode)
	}

	encoded, stdErr := json.Marshal(data)
	if stdErr != nil {
		return "", ErrWrapperInvalidData.Wrap(stdErr).AddCategory(ErrTypeEncode)
	}

	return string(encoded), nil
}
