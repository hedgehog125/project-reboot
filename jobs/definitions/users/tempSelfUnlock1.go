package users

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
)

type TempSelfUnlock1Body struct {
	Username string `binding:"required" json:"username"`
}

func TempSelfUnlock1(app *common.App) *jobs.Definition {
	return &jobs.Definition{
		ID:            "TEMP_SELF_UNLOCK",
		Version:       1,
		Weight:        1,
		NoParallelize: true,
		BodyType:      &TempSelfUnlock1Body{},
		Handler: func(jobCtx *jobs.Context) error {
			body := &TempSelfLock1Body{}
			jobErr := jobCtx.Decode(body)
			if jobErr != nil {
				return jobErr
			}

			return dbcommon.WithWriteTx(
				jobCtx.Context, app.Database,
				func(tx *ent.Tx, ctx context.Context) error {
					userOb, stdErr := tx.User.Query().
						Where(user.Username(body.Username)).
						Only(ctx)
					if stdErr != nil {
						return stdErr
					}
					if userOb.LockedUntil == nil {
						jobCtx.Logger.Info(
							"didn't need to unlock the user because they are already unlocked",
							"userID", userOb.ID,
						)
						return nil
					}
					userOb, stdErr = userOb.Update().
						ClearLockedUntil().
						Save(ctx)
					if stdErr != nil {
						return stdErr
					}
					_, stdErr = tx.Session.Delete().
						Where(session.UserID(userOb.ID)).
						Exec(ctx)
					if stdErr != nil {
						return stdErr
					}

					_, _, commErr := app.Messengers.SendUsingAll(
						&common.Message{
							Type: common.MessageSelfUnlock,
							User: userOb,
							Time: body.Until,
						},
						ctx,
					)
					if commErr != nil {
						return commErr
					}
					jobCtx.Logger.Info(
						"user was unlocked because the self-lock expired",
						"userID", userOb.ID,
						"isLocked", userOb.Locked,
					)
					return nil
				},
			)
		},
	}
}
