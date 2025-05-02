package twofactoractions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hedgehog125/project-reboot/ent"
)

var ErrUnknownActionType = errors.New("unknown action type")
var ErrInvalidData = errors.New("invalid action data")

type actionFuncType = func(action *ent.TwoFactorAction) error

var actionMap = map[string]actionFuncType{
	"TEMP_SELF_LOCK_1": func(action *ent.TwoFactorAction) error {
		type actionDataType struct {
			Username string `binding:"required" json:"username"`
		}

		parsed := actionDataType{}
		err := json.Unmarshal([]byte(action.Data), &parsed)
		if err != nil {
			// TODO: add the JSON decode error to this
			return ErrInvalidData
		}

		return nil
	},
}

func Confirm(action *ent.TwoFactorAction) error {
	// TODO: move expiry checking to this

	fullID := fmt.Sprintf("%v_%v", action.Type, action.Version)
	actionFunc, ok := actionMap[fullID]
	if !ok {
		return ErrUnknownActionType
	}

	return actionFunc(action)
}
