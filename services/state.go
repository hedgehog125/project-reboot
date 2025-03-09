package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/core"
)

func InitState() *common.State {
	state := common.State{
		AdminCode: common.InitChannel([]byte{}),
	}

	core.UpdateAdminCode(&state)

	return &state
}
