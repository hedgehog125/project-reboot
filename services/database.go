package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/ent"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(env *common.Env) *Database {
	return &Database{
		env:       env,
		readyChan: make(chan struct{}),
	}
}

type Database struct {
	env       *common.Env
	client    *ent.Client
	readyChan chan struct{}
}

func (service *Database) Start() {
	defer close(service.readyChan)

	err := os.MkdirAll(service.env.MOUNT_PATH, 0700)
	if err != nil {
		log.Fatalf("couldn't create storage directory. Error:\n%v", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%v?_fk=1&_busy_timeout=10000", filepath.Join(service.env.MOUNT_PATH, "database.sqlite3")))
	if err != nil {
		log.Fatalf("couldn't start database. Error:\n%v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	driver := ent.Driver(entsql.OpenDB("sqlite3", db))
	client := ent.NewClient(driver)
	service.client = client

	// TODO: does not cancelling this cause a memory leak?
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Schema.Create(ctx)
	if err != nil {
		_ = client.Close()
		log.Fatalf("couldn't create schema resources. Error:\n%v", err)
	}
}

func (service *Database) Client() *ent.Client {
	<-service.readyChan
	return service.client
}
func (service *Database) ReadTx(ctx context.Context) (*ent.Tx, error) {
	<-service.readyChan
	return service.client.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
}
func (service *Database) WriteTx(ctx context.Context) (*ent.Tx, error) {
	<-service.readyChan
	return service.client.Tx(ctx)
}

func (service *Database) Shutdown() {
	<-service.readyChan
	err := service.client.Close()
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the database:\n%v\n", err.Error())
	}
}
