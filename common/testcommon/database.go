package testcommon

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/hedgehog125/project-reboot/ent"
)

type TestDatabase struct {
	client *ent.Client
}

var dbCounter atomic.Int64

func CreateDB() *TestDatabase {
	// TODO: review options
	// TODO: what does shared cache do any why is it sometimes necessary in order to stop the database being deleted mid test?
	// ^ this seems to enable WAL mode? Which isn't what I want
	db, stdErr := sql.Open("sqlite3", fmt.Sprintf("file:temp%v?mode=memory&cache=shared&_fk=1&_busy_timeout=10000&_txlock=immediate", dbCounter.Add(1)))
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
		client: client,
	}
}
func (db *TestDatabase) Start() {
	// TODO: move initialisation logic into here like the real DB service?
}
func (db *TestDatabase) Client() *ent.Client {
	return db.client
}
func (db *TestDatabase) ReadTx(ctx context.Context) (*ent.Tx, error) {
	return db.client.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
}
func (db *TestDatabase) WriteTx(ctx context.Context) (*ent.Tx, error) {
	return db.client.Tx(ctx)
}
func (db *TestDatabase) Shutdown() {
	err := db.client.Close()
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the database:\n%v\n", err.Error())
	}
}
