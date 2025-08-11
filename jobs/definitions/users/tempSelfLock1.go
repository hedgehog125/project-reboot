package users

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
)

// TODO: these types need to go somewhere else so that the services package can run jobs? Maybe not since the jobs package doesn't actually depend on too much?
type TempSelfLock1Body struct {
	Username string               `binding:"required" json:"username"`
	Until    common.ISOTimeString `binding:"required" json:"until"`
}

func TempSelfLock1(app *common.App) *jobs.Definition {
	return &jobs.Definition{
		ID:            "TEMP_SELF_LOCK",
		Version:       1,
		Priority:      jobs.HighPriority,
		Weight:        1,
		NoParallelize: true,
		BodyType:      &TempSelfLock1Body{},
		Handler: func(jobCtx *jobs.Context) error {
			body := &TempSelfLock1Body{}
			jobErr := jobCtx.Decode(body)
			if jobErr != nil {
				return jobErr
			}

			return dbcommon.WithWriteTx(jobCtx.Context, app.Database, func(tx *ent.Tx, ctx context.Context) error {
				userOb, stdErr := tx.User.Query().
					Where(user.Username(body.Username)).
					Only(ctx)
				if stdErr != nil {
					return common.ErrWrapperDatabase.Wrap(stdErr)
				}
				userOb, stdErr = userOb.Update().SetLockedUntil(body.Until.Time).
					Save(ctx)
				if stdErr != nil {
					return common.ErrWrapperDatabase.Wrap(stdErr)
				}

				return app.Messengers.SendUsingAll(
					&common.Message{
						Type:  common.MessageSelfLock,
						User:  userOb,
						Until: body.Until.Time,
					},
					ctx,
				).StandardError()
			})
		},
	}
}
