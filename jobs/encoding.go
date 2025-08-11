package jobs

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hedgehog125/project-reboot/common"
)

func (registry *Registry) Encode(versionedType string, data any) (string, *common.Error) {
	actionDef, ok := registry.jobs[versionedType]
	if !ok {
		return "", ErrWrapperEncode.Wrap(ErrUnknownJobType)
	}

	dataType := reflect.TypeOf(data)
	if dataType != actionDef.reflectedBodyType {
		return "", ErrWrapperEncode.Wrap(ErrWrapperInvalidData.Wrap(
			fmt.Errorf("data type %s isn't the expected type %s",
				dataType, actionDef.reflectedBodyType),
		))
	}

	encoded, stdErr := json.Marshal(data)
	if stdErr != nil {
		return "", ErrWrapperEncode.Wrap(ErrWrapperInvalidData.Wrap(stdErr))
	}

	return string(encoded), nil
}
