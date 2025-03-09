package common

import (
	"context"

	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
)

func ReadMessageUserInfo(username string, dbClient *ent.Client) (*MessageUserInfo, error) {
	row, err := dbClient.User.Query().
		Where(user.Username(username)).
		Select(user.FieldAlertDiscordId, user.FieldAlertEmail).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	return &MessageUserInfo{
		Username:       username,
		AlertDiscordId: row.AlertDiscordId,
		AlertEmail:     row.AlertEmail,
	}, nil
}
