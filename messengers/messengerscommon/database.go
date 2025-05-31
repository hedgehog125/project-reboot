package messengerscommon

import (
	"context"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
)

func ReadMessageUserInfo(username string, dbClient *ent.Client) (*common.MessageUserInfo, *common.Error) {
	row, err := dbClient.User.Query().
		Where(user.Username(username)).
		Select(user.FieldAlertDiscordId, user.FieldAlertEmail).
		Only(context.Background())
	if err != nil {
		return nil, ErrWrapperDatabase.Wrap(err)
	}

	//exhaustruct:enforce
	return &common.MessageUserInfo{
		Username:       username,
		AlertDiscordId: row.AlertDiscordId,
		AlertEmail:     row.AlertEmail,
	}, nil
}
