package users

import (
	"context"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/session"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/jobs"
)

type TempSelfLock1Body struct {
	Username string    `binding:"required" json:"username"`
	Until    time.Time `binding:"required" json:"until"`
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
					return stdErr
				}
				userOb, stdErr = userOb.Update().SetLockedUntil(body.Until).
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
						Type: common.MessageSelfLock,
						User: userOb,
						Time: body.Until,
					},
					ctx,
				)
				if commErr != nil {
					return commErr
				}
				jobCtx.Logger.Info(
					"user has successfully self-locked",
					"userID", userOb.ID,
				)
				return nil
			})
		},
	}
}
