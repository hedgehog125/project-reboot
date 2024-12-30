package subfns

import (
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() *intertypes.Env {
	_ = godotenv.Load(".env")

	env := intertypes.Env{
		MOUNT_PATH: util.RequireEnv("MOUNT_PATH"),
	}
	return &env
}
