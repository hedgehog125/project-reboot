package twofactoractions

import (
	"encoding/json"
	"time"
)

type NoAction1 struct {
	Foo string
}
type TempSelfLock1 struct {
	Username string    `binding:"required" json:"username"`
	Until    time.Time `binding:"required" json:"until"`
}

// TODO: combine into single map
var actionTypeMap = map[string]any{
	"No_ACTION_1":      NoAction1{}, // For testing
	"TEMP_SELF_LOCK_1": TempSelfLock1{},
}

// TODO: pass registry as argument for better testing?
func Encode(fullType string, data any) (string, error) {
	actionType, ok := actionTypeMap[fullType]
	if !ok {
		return "", ErrUnknownActionType
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", ErrInvalidData
	}

	// TODO: is there a better way to do this? With reflection maybe?
	err = json.Unmarshal(encoded, &actionType)
	if err != nil {
		return "", ErrInvalidData
	}

	return string(encoded), nil
}
