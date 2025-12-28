package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/ent"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(app *common.App) *Database {
	return &Database{
		app: app,
	}
}

type Database struct {
	app          *common.App
	client       *ent.Client
	startOnce    sync.Once
	shutdownOnce sync.Once
	mu           sync.RWMutex
}

func (service *Database) Start() {
	service.startOnce.Do(func() {
		stdErr := os.MkdirAll(service.app.Env.MOUNT_PATH, 0700)
		if stdErr != nil {
			log.Fatalf("couldn't create storage directory. error:\n%v", stdErr)
		}

		db, stdErr := sql.Open("sqlite3", fmt.Sprintf(
			"%v?_fk=1&_foreign_keys=on",
			filepath.Join(service.app.Env.MOUNT_PATH, "database.sqlite3"),
		))
		if stdErr != nil {
			log.Fatalf("couldn't start database. error:\n%v", stdErr)
		}

		db.SetMaxIdleConns(5)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
		driver := ent.Driver(entsql.OpenDB("sqlite3", db))
		client := ent.NewClient(driver)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stdErr = client.Schema.Create(ctx)
		if stdErr == nil {
			service.mu.Lock()
			service.client = client
			service.mu.Unlock()
		} else {
			_ = client.Close()
			log.Fatalf("couldn't create schema resources. error:\n%v", stdErr)
		}
	})
}

func (service *Database) Client() *ent.Client {
	service.mu.RLock()
	defer service.mu.RUnlock()
	if service.client == nil {
		panic("can't get database client, service isn't running")
	}
	return service.client
}
func (service *Database) ReadTx(ctx context.Context) (*ent.Tx, error) {
	return service.Client().BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
}
func (service *Database) WriteTx(ctx context.Context) (*ent.Tx, error) {
	return service.Client().Tx(ctx)
}

func (service *Database) Shutdown() {
	service.shutdownOnce.Do(func() {
		service.mu.Lock()
		client := service.client
		service.client = nil
		service.mu.Unlock()

		if client != nil {
			stdErr := client.Close()
			if stdErr != nil {
				service.app.Logger.Warn("an error occurred while shutting down the database", stdErr)
			}
		}
	})
}
func (service *Database) DefaultLogger() common.Logger {
	return service.app.Logger
}
