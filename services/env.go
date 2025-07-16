package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() *common.Env {
	_ = godotenv.Load(".env") // TODO: log some errors

	env := common.Env{
		IS_DEV:                        common.RequireBoolEnv("IS_DEV"),
		PORT:                          common.RequireIntEnv("PORT"),
		MOUNT_PATH:                    common.RequireEnv("MOUNT_PATH"),
		PROXY_ORIGINAL_IP_HEADER_NAME: common.RequireEnv("PROXY_ORIGINAL_IP_HEADER_NAME"),

		JOB_POLL_INTERVAL:    common.RequireSecondsEnv("JOB_POLL_INTERVAL"),
		MAX_TOTAL_JOB_WEIGHT: common.RequireIntEnv("MAX_TOTAL_JOB_WEIGHT"),

		UNLOCK_TIME:         common.RequireSecondsEnv("UNLOCK_TIME"),
		AUTH_CODE_VALID_FOR: common.RequireSecondsEnv("AUTH_CODE_VALID_FOR"),

		PASSWORD_HASH_SETTINGS: &common.PasswordHashSettings{
			Time:    common.RequireUint32Env("PASSWORD_HASH_TIME"),
			Memory:  common.RequireUint32Env("PASSWORD_HASH_MEMORY"),
			Threads: common.RequireUint8Env("PASSWORD_HASH_THREADS"),
		},

		DISCORD_TOKEN:  common.OptionalEnv("DISCORD_TOKEN", ""),
		SENDGRID_TOKEN: common.OptionalEnv("SENDGRID_TOKEN", ""),
	}
	return &env
}
