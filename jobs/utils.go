package jobs

// Utils for job definitions to use

// TODO: how should the errors work with this?
// func TxHandler(
// 	dbClient *ent.Client,
// 	handler func(ctx *Context, tx *ent.Tx) error,
// ) HandlerFunc {
// 	return func(ctx *Context) error {
// 		return dbcommon.WithTx(
// 			ctx.Context,
// 			dbClient,
// 			func(tx *ent.Tx) error {
// 				return handler(ctx, tx)
// 			},
// 		)
// 	}
// }
