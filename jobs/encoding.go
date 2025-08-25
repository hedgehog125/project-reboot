package jobs

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hedgehog125/project-reboot/common"
)

func (registry *Registry) Encode(versionedType string, body any) (json.RawMessage, *common.Error) {
	actionDef, ok := registry.jobs[versionedType]
	if !ok {
		return nil, ErrWrapperEncode.Wrap(ErrUnknownJobType)
	}

	bodyType := reflect.TypeOf(body)
	if bodyType != actionDef.reflectedBodyType {
		return nil, ErrWrapperEncode.Wrap(ErrWrapperInvalidBody.Wrap(
			fmt.Errorf("body type %s isn't the expected type %s",
				bodyType, actionDef.reflectedBodyType),
		))
	}

	encoded, stdErr := json.Marshal(body)
	if stdErr != nil {
		return nil, ErrWrapperEncode.Wrap(ErrWrapperInvalidBody.Wrap(stdErr))
	}

	return encoded, nil
}
