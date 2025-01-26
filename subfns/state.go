package subfns

import (
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
)

func InitState() *intertypes.State {
	state := intertypes.State{
		AdminCode: util.InitChannel([]byte{}),
	}

	core.UpdateAdminCode(&state)

	return &state
}
