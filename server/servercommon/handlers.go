package servercommon

import (
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

func WithTx(
	app *ServerApp,
	handler func(ctx *gin.Context, tx *ent.Tx) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := dbcommon.WithTx(
			ctx,
			app.Database,
			func(tx *ent.Tx) error {
				return handler(ctx, tx)
			},
		)
		if err != nil {
			ctx.Error(err)
		}
	}
}
