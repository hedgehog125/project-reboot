package users

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
	"github.com/hedgehog125/project-reboot/messengers/messengerscommon"
)

// TODO: these types need to go somewhere else so that the services package can run jobs? Maybe not since the jobs package doesn't actually depend on too much?
type TempSelfLock1Body struct {
	Username string               `binding:"required" json:"username"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}

func TempSelfLock1(app *common.App) *jobs.Definition {
	return &jobs.Definition{
		ID:       "TEMP_SELF_LOCK",
		Version:  1,
		Priority: jobs.HighPriority,
		BodyType: &TempSelfLock1Body{},
		Handler: func(jobCtx *jobs.Context) error {
			body := &TempSelfLock1Body{}
			jobErr := jobCtx.Decode(body)
			if jobErr != nil {
				return jobErr
			}

			return dbcommon.WithWriteTx(jobCtx.Context, app.Database, func(tx *ent.Tx, ctx context.Context) error {
				_, stdErr := tx.User.Update().
					Where(user.Username(body.Username)).
					SetLockedUntil(body.Until.Time).Save(ctx)
				if stdErr != nil {
					return common.ErrWrapperDatabase.Wrap(stdErr)
				}

				userInfo, commErr := messengerscommon.ReadUserContacts(body.Username, ctx)
				if commErr != nil {
					return commErr
				}

				commErr = app.Messengers.SendUsingAll(
					common.Message{
						Type:  common.MessageSelfLock,
						User:  userInfo,
						Until: body.Until.Time,
					},
					ctx,
				)
				if commErr != nil {
					return commErr
				}

				return nil
			})
		},
	}
}
