package jobs

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

// Utils for job definitions to use

func TxHandler(
	db common.DatabaseService,
	handler func(ctx *Context, tx *ent.Tx) error,
) HandlerFunc {
	return func(ctx *Context) error {
		return dbcommon.WithTx(
			ctx.Context,
			db,
			func(tx *ent.Tx) error {
				return handler(ctx, tx)
			},
		)
	}
}
