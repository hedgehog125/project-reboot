package services

import (
	"context"
	"log"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/NicoClack/cryptic-stash/backend/keyvalue"
	"github.com/NicoClack/cryptic-stash/backend/keyvalue/definitions"
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
			return service.InitAll(ctx)
		},
	)
	if stdErr != nil {
		log.Fatalf("failed to initialize key/value service. error:\n%v", stdErr.Error())
	}
}
