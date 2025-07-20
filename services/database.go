package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(env *common.Env) common.DatabaseService {
	return &databaseService{
		env:       env,
		readyChan: make(chan struct{}),
	}
}

type databaseService struct {
	env       *common.Env
	client    *ent.Client
	readyChan chan struct{}
}

func (service *databaseService) Start() {
	defer close(service.readyChan)

	err := os.MkdirAll(service.env.MOUNT_PATH, 0700)
	if err != nil {
		log.Fatalf("couldn't create storage directory. Error:\n%v", err)
	}

	client, err := ent.Open("sqlite3", fmt.Sprintf("%v?&_fk=1", filepath.Join(service.env.MOUNT_PATH, "database.sqlite3")))
	if err != nil {
		log.Fatalf("couldn't start database. Error:\n%v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Schema.Create(ctx)
	if err != nil {
		_ = client.Close()
		log.Fatalf("couldn't create schema resources. Error:\n%v", err)
	}
}

func (service *databaseService) Client() *ent.Client {
	<-service.readyChan
	return service.client
}
func (service *databaseService) Tx(ctx context.Context) (*ent.Tx, error) {
	<-service.readyChan
	return service.client.Tx(ctx)
}

func (service *databaseService) Shutdown() {
	<-service.readyChan
	err := service.client.Close()
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the database:\n%v\n", err.Error())
	}
}
