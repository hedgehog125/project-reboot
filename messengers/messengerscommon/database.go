package messengerscommon

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
)

func ReadUserContacts(username string, ctx context.Context) (*common.UserContacts, *common.Error) {
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return nil, ErrNoTxInContext.AddCategory(ErrTypeReadUserContacts)
	}
	row, err := tx.User.Query().
		Where(user.Username(username)).
		Select(user.FieldAlertDiscordId, user.FieldAlertEmail).
		Only(ctx)
	if err != nil {
		return nil, ErrWrapperDatabase.Wrap(err).AddCategory(ErrTypeReadUserContacts)
	}

	//exhaustruct:enforce
	return &common.UserContacts{
		Username:       username,
		AlertDiscordId: row.AlertDiscordId,
		AlertEmail:     row.AlertEmail,
	}, nil
}
