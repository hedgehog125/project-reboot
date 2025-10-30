package services

import (
	"context"
	"log"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/keyvalue"
	"github.com/hedgehog125/project-reboot/keyvalue/definitions"
)

type KeyValue struct {
	*keyvalue.Registry
}

func NewKeyValue(app *common.App, registerFuncs ...func(group *keyvalue.RegistryGroup)) *KeyValue {
	registry := keyvalue.NewRegistry(app)
	definitions.Register(registry.Group(""))
	for _, registerFunc := range registerFuncs {
		registerFunc(registry.Group(""))
	}

	return &KeyValue{
		Registry: registry,
	}
}

func (service *KeyValue) Init() {
	stdErr := dbcommon.WithWriteTx(
		context.TODO(), service.App.Database,
		func(tx *ent.Tx, ctx context.Context) error {
			// TODO: why no transaction?
			return service.Registry.InitAll(ctx).StandardError()
		},
	)
	if stdErr != nil {
		log.Fatalf("failed to initialize key/value service. error:\n%v", stdErr.Error())
	}
}
