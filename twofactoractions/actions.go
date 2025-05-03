package twofactoractions

import (
	"encoding/json"

	"github.com/hedgehog125/project-reboot/ent"
)

type actionFuncType = func(action *ent.TwoFactorAction) error

var actionMap = map[string]actionFuncType{
	"TEMP_SELF_LOCK_1": func(action *ent.TwoFactorAction) error {
		parsed := TempSelfLock1{}
		err := json.Unmarshal([]byte(action.Data), &parsed)
		if err != nil {
			// TODO: add the JSON decode error to this
			return ErrInvalidData
		}

		return nil
	},
}
