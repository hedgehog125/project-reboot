package testcommon

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/ent"
)

type TestDatabase struct {
	client       *ent.Client
	logger       common.Logger
	startTxHooks []func(tx *ent.Tx) error
}

var (
	dbCounter = int64(0)
	// TODO: this seems to be necessary because of some race conditions in Ent/Atlas
	createMu = sync.Mutex{}
)

func CreateDB() *TestDatabase {
	// TODO: review options
	// TODO: what does shared cache do any why is it sometimes necessary
	// in order to stop the database being deleted mid test?
	// ^ this seems to enable WAL mode? Which isn't what I want
	createMu.Lock()
	defer createMu.Unlock()
	dbCounter++
	db, stdErr := sql.Open("sqlite3", fmt.Sprintf(
		"file:temp%v?mode=memory&cache=shared&_fk=1&_busy_timeout=250&_foreign_keys=on",
		dbCounter,
	))
	if stdErr != nil {
		panic(fmt.Sprintf("failed to open test database. error: %v", stdErr.Error()))
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	driver := ent.Driver(entsql.OpenDB("sqlite3", db))
	client := ent.NewClient(driver)

	stdErr = client.Schema.Create(context.Background())
	if stdErr != nil {
		client.Close()
		panic(fmt.Sprintf("failed to create test database schema. error: %v", stdErr.Error()))
	}

	return &TestDatabase{
		client:       client,
		logger:       common.GetLogger(context.Background(), nil),
		startTxHooks: []func(tx *ent.Tx) error{},
	}
}
func (db *TestDatabase) Start() {
	// TODO: move initialisation logic into here like the real DB service?
}
func (db *TestDatabase) Client() *ent.Client {
	return db.client
}
func (db *TestDatabase) ReadTx(ctx context.Context) (*ent.Tx, error) {
	tx, stdErr := db.client.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if stdErr != nil {
		return nil, stdErr
	}
	stdErr = db.runStartTxHooks(tx)
	if stdErr != nil {
		_ = tx.Rollback()
		return nil, stdErr
	}
	return tx, nil
}
func (db *TestDatabase) WriteTx(ctx context.Context) (*ent.Tx, error) {
	tx, stdErr := db.client.Tx(ctx)
	if stdErr != nil {
		return nil, stdErr
	}
	stdErr = db.runStartTxHooks(tx)
	if stdErr != nil {
		_ = tx.Rollback()
		return nil, stdErr
	}
	return tx, nil
}
func (db *TestDatabase) Shutdown() {
	stdErr := db.client.Close()
	if stdErr != nil {
		db.logger.Warn("an error occurred while shutting down a test database", "error", stdErr)
	}
}
func (db *TestDatabase) DefaultLogger() common.Logger {
	return db.logger
}

func (db *TestDatabase) AddStartTxHook(hook func(tx *ent.Tx) error) {
	db.startTxHooks = append(db.startTxHooks, hook)
}
func (db *TestDatabase) runStartTxHooks(tx *ent.Tx) error {
	for _, hook := range db.startTxHooks {
		stdErr := hook(tx)
		if stdErr != nil {
			return stdErr
		}
	}
	return nil
}
