package testcommon

import (
	"context"
	"fmt"
	"sync"

	"github.com/hedgehog125/project-reboot/ent"
)

type TestDatabase struct {
	client *ent.Client
}

func CreateDB() *TestDatabase {
	// TODO: review options
	client, stdErr := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if stdErr != nil {
		panic(fmt.Sprintf("failed to open test database. error: %v", stdErr.Error()))
	}

	// Run the auto migration tool.
	// This will create the necessary tables based on your Ent schemas.
	stdErr = client.Schema.Create(context.Background())
	if stdErr != nil {
		client.Close()
		panic(fmt.Sprintf("failed to create test database schema. error: %v", stdErr.Error()))
	}

	return &TestDatabase{
		client: client,
	}
}
func (db *TestDatabase) Client() *ent.Client {
	return db.client
}
func (db *TestDatabase) TwoFactorActionMutex() *sync.Mutex { // TODO: remove
	panic("not implemented")
}
func (db *TestDatabase) Shutdown() {
	db.client.Close()
}
