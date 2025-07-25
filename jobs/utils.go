package jobs

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

// Utils for job definitions to use

func ReadTxHandler(
	db common.DatabaseService,
	handler func(ctx *Context, tx *ent.Tx) error,
) HandlerFunc {
	return func(jobCtx *Context) error {
		return dbcommon.WithReadTx(
			jobCtx.Context,
			db,
			func(tx *ent.Tx, ctx context.Context) error {
				jobCtx.Context = ctx // TODO: is this right?
				return handler(jobCtx, tx)
			},
		)
	}
}
