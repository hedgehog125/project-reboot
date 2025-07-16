package jobs

// Utils for job definitions to use

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

func TxHandler(
	dbClient *ent.Client,
	handler func(ctx *Context, tx *ent.Tx) *common.Error,
) HandlerFunc {
	return func(ctx *Context) *common.Error {
		return dbcommon.WithTx(
			ctx.Context,
			dbClient,
			func(tx *ent.Tx) *common.Error {
				return handler(ctx, tx)
			},
		)
	}
}
