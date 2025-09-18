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

func NewDatabase(app *common.App) *Database {
	return &Database{
		app:       app,
		readyChan: make(chan struct{}),
	}
}

type Database struct {
	app       *common.App
	client    *ent.Client
	readyChan chan struct{}
}

func (service *Database) Start() {
	defer close(service.readyChan)

	err := os.MkdirAll(service.app.Env.MOUNT_PATH, 0700)
	if err != nil {
		log.Fatalf("couldn't create storage directory. Error:\n%v", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%v?_fk=1&_busy_timeout=250&_foreign_keys=on", filepath.Join(service.app.Env.MOUNT_PATH, "database.sqlite3")))
	if err != nil {
		log.Fatalf("couldn't start database. error:\n%v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	driver := ent.Driver(entsql.OpenDB("sqlite3", db))
	client := ent.NewClient(driver)
	service.client = client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Schema.Create(ctx)
	if err != nil {
		_ = client.Close()
		log.Fatalf("couldn't create schema resources. error:\n%v", err)
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
	stdErr := service.client.Close()
	if stdErr != nil {
		service.app.Logger.Warn("an error occurred while shutting down the database", stdErr)
	}
}
func (service *Database) DefaultLogger() common.Logger {
	return service.app.Logger
}
