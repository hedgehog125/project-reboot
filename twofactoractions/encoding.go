package twofactoractions

import (
	"encoding/json"
)

type NoAction1 struct {
	Foo string
}

func (registry *Registry) Encode(fullType string, data any) (string, error) {
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return "", ErrUnknownActionType
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", ErrInvalidData
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType
	err = json.Unmarshal(encoded, &temp)
	if err != nil {
		return "", ErrInvalidData
	}

	return string(encoded), nil
}
