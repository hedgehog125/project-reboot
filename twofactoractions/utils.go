package twofactoractions

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
)

func TxHandler(
	dbClient *ent.Client,
	handler func(action *Action, tx *ent.Tx) *common.Error,
) HandlerFunc {
	return func(action *Action) *common.Error {
		return dbcommon.WithTx(
			action.Context,
			dbClient,
			func(tx *ent.Tx) *common.Error {
				return handler(action, tx)
			},
		)
	}
}
