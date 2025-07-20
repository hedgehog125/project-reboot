package users

import (
	"fmt"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
)

type TempSelfLock1Body struct {
	Username string               `binding:"required" json:"username"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}

func TempSelfLock1(app *common.App) *jobs.Definition {
	dbClient := app.Database.Client()

	return &jobs.Definition{
		ID:       "TEMP_SELF_LOCK",
		Version:  1,
		Priority: jobs.HighPriority,
		BodyType: &TempSelfLock1Body{},
		Handler: jobs.TxHandler(dbClient, func(ctx *jobs.Context, tx *ent.Tx) error {
			body := &TempSelfLock1Body{}
			jobErr := ctx.Decode(body)
			if jobErr != nil {
				return jobErr
			}

			_, stdErr := tx.User.Update().
				Where(user.Username(body.Username)).
				SetLockedUntil(body.Until.Time).Save(ctx.Context)
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
