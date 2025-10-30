package services

import "github.com/hedgehog125/project-reboot/common"

type KeyValue struct {
	App *common.App
}

func NewKeyValue(app *common.App) *KeyValue {
	return &KeyValue{
		App: app,
	}
}
