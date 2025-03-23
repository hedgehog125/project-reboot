package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(env *common.Env) common.DatabaseService {
	err := os.MkdirAll(env.MOUNT_PATH, 0700)
	if err != nil {
		log.Fatalf("couldn't create storage directory. Error:\n%v", err)
	}

	client, err := ent.Open("sqlite3", fmt.Sprintf("%v?&_fk=1", filepath.Join(env.MOUNT_PATH, "database.sqlite3")))
	if err != nil {
		log.Fatalf("couldn't start database. Error:\n%v", err)
	}

	err = client.Schema.Create(context.Background())
	if err != nil {
		client.Close()
		log.Fatalf("couldn't create schema resources. Error:\n%v", err)
	}

	return &databaseService{
		client: client,
	}
}

type databaseService struct {
	client *ent.Client
}

func (service *databaseService) Client() *ent.Client {
	return service.client
}

func (service *databaseService) Shutdown() {
	err := service.client.Close()
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the database:\n%v\n", err.Error())
	}
}

func (service *databaseService) ReadMessageUserInfo(username string) (*common.MessageUserInfo, error) {
	row, err := service.client.User.Query().
		Where(user.Username(username)).
		Select(user.FieldAlertDiscordId, user.FieldAlertEmail).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	//exhaustruct:enforce
	return &common.MessageUserInfo{
		Username:       username,
		AlertDiscordId: row.AlertDiscordId,
		AlertEmail:     row.AlertEmail,
	}, nil
}
