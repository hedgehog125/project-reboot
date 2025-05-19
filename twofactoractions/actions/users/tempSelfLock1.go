package users

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type TempSelfLock1Body struct {
	Username string               `binding:"required" json:"username"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}

func TempSelfLock1(app *common.App) twofactoractions.ActionDefinition[any] {
	return twofactoractions.ActionDefinition[any]{
		ID:      "TEMP_SELF_LOCK",
		Version: 1,
		Handler: func(action *twofactoractions.Action[any]) error {
			// TODO: do I really have to cast it back to the struct type? Can I improve the generics?

			// TODO: implement
			return nil
		},
		BodyType: TempSelfLock1Body{},
	}
}
