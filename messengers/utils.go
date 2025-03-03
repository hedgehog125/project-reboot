package messengers

import (
	"context"

	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
)

func ReadUserInfo(username string, dbClient *ent.Client) (*UserInfo, error) {
	row, err := dbClient.User.Query().
		Where(user.Username(username)).
		Select(user.FieldAlertDiscordId, user.FieldAlertEmail).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Username:       username,
		AlertDiscordId: row.AlertDiscordId,
		AlertEmail:     row.AlertEmail,
	}, nil
}
