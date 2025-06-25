package users

import (
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

func TempSelfLock1(app *common.App) twofactoractions.ActionDefinition[any] {
	dbClient := app.Database.Client()

	return twofactoractions.ActionDefinition[any]{
		ID:      "TEMP_SELF_LOCK",
		Version: 1,
		Handler: twofactoractions.TxHandler(dbClient, func(action *twofactoractions.Action[any], tx *ent.Tx) *common.Error {
			// TODO: this is getting turned into a map somehow and not casting back?
			body, ok := action.Body.(*TempSelfLock1Body)
			if !ok {
				return twofactoractions.ErrActionInvalidBody.Clone()
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

			app.Messenger.SendUsingAll(common.Message{
				Type: common.MessageSelfLock,
				User: userInfo,
			})

			return nil
		}),
		BodyType: TempSelfLock1Body{},
	}
}
