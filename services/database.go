package services

import (
	"context"
	"database/sql"
	"errors"
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

var ErrDatabaseNotStarted = errors.New("database service not started")
var ErrDatabaseStartFailed = errors.New("database service failed to start")
var ErrDatabaseAlreadyShutdown = errors.New("database service has already been shut down")

func NewDatabase(app *common.App) *Database {
	return &Database{
		app:          app,
		readyChan:    make(chan struct{}),
		startingChan: make(chan struct{}),
	}
}

type Database struct {
	app                 *common.App
	client              *ent.Client
	startingChan        chan struct{}
	readyChan           chan struct{}
	startFailed         bool
	shutdownOnce        sync.Once
	isShutdownCompleted bool
	mu                  sync.RWMutex
}

func (service *Database) Start() {
	if service.markAsStarting() {
		<-service.readyChan
		return
	}
	service.assertNotShutdown()

	stdErr := os.MkdirAll(service.app.Env.MOUNT_PATH, 0700)
	if stdErr != nil {
		log.Fatalf("couldn't create storage directory. Error:\n%v", stdErr)
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
	service.client = client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stdErr = client.Schema.Create(ctx)
	if stdErr != nil {
		service.mu.Lock()
		service.startFailed = true
		service.mu.Unlock()
		_ = client.Close()
		log.Fatalf("couldn't create schema resources. error:\n%v", stdErr)
	}
	close(service.readyChan)
}
func (service *Database) markAsStarting() bool {
	service.mu.Lock()
	defer service.mu.Unlock()
	select {
	case <-service.startingChan:
		return true
	default:
		close(service.startingChan)
	}
	return false
}
func (service *Database) assertNotShutdown() {
	service.mu.RLock()
	defer service.mu.RUnlock()
	if service.isShutdownCompleted {
		panic(ErrDatabaseAlreadyShutdown)
	}
}
func (service *Database) waitForStart() error {
	select {
	case <-service.startingChan:
	default:
		return ErrDatabaseNotStarted
	}
	<-service.readyChan

	service.mu.RLock()
	defer service.mu.RUnlock()
	if service.startFailed {
		return ErrDatabaseStartFailed
	}
	return nil
}

func (service *Database) Client() *ent.Client {
	common.PanicIfError(service.waitForStart())
	return service.client
}
func (service *Database) ReadTx(ctx context.Context) (*ent.Tx, error) {
	stdErr := service.waitForStart()
	if stdErr != nil {
		return nil, stdErr
	}
	return service.client.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
}
func (service *Database) WriteTx(ctx context.Context) (*ent.Tx, error) {
	stdErr := service.waitForStart()
	if stdErr != nil {
		return nil, stdErr
	}
	return service.client.Tx(ctx)
}

func (service *Database) Shutdown() {
	service.shutdownOnce.Do(func() {
		defer func() {
			service.mu.Lock()
			service.isShutdownCompleted = true
			service.mu.Unlock()
		}()

		select {
		case <-service.startingChan:
			<-service.readyChan
			service.mu.RLock()
			startFailed := service.startFailed
			service.mu.RUnlock()
			if startFailed {
				return
			}

			stdErr := service.client.Close()
			if stdErr != nil {
				service.app.Logger.Warn("an error occurred while shutting down the database", stdErr)
			}
		default:
			return
		}
	})
}
func (service *Database) DefaultLogger() common.Logger {
	return service.app.Logger
}
