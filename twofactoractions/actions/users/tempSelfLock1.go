package users

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type TempSelfLock1Body struct {
	Username string               `binding:"required" json:"username"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}

func TempSelfLock1(app *common.App) twofactoractions.ActionDefinition {
	dbClient := app.Database.Client()

	return twofactoractions.ActionDefinition{
		ID:       "TEMP_SELF_LOCK",
		Version:  1,
		BodyType: func() any { return &TempSelfLock1Body{} },
		Handler: twofactoractions.TxHandler(dbClient, func(action *twofactoractions.Action, tx *ent.Tx) *common.Error {
			body := &TempSelfLock1Body{}
			commErr := action.Decode(body)
			if commErr != nil {
				return commErr
			}

			_, stdErr := tx.User.Update().
				Where(user.Username(body.Username)).
				SetLockedUntil(body.Until.Time).Save(action.Context)
			if stdErr != nil {
				return common.ErrWrapperDatabase.Wrap(stdErr)
			}

			userInfo, commErr := messengerscommon.ReadMessageUserInfo(body.Username, dbClient)
			if commErr != nil {
				return commErr
			}

			errs := app.Messenger.SendUsingAll(common.Message{
				Type:  common.MessageSelfLock,
				User:  userInfo,
				Until: body.Until.Time,
			})
			fmt.Println(errs) // TODO

			return nil
		}),
	}
}
