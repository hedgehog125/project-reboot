package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/intertypes"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase(env *intertypes.Env) *ent.Client {
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

	return client
}

func ShutdownDatabase(client *ent.Client) {
	err := client.Close()
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the database:\n%v\n", err.Error())
	}
}
