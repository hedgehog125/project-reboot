package services

import (
	"github.com/hedgehog125/project-reboot/common"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() *common.Env {
	_ = godotenv.Load(".env")

	env := common.Env{
		IS_DEV:                        common.RequireBoolEnv("IS_DEV"),
		PORT:                          common.RequireIntEnv("PORT"),
		MOUNT_PATH:                    common.RequireEnv("MOUNT_PATH"),
		PROXY_ORIGINAL_IP_HEADER_NAME: common.RequireEnv("PROXY_ORIGINAL_IP_HEADER_NAME"),
		UNLOCK_TIME:                   common.RequireInt64Env("UNLOCK_TIME"),
		AUTH_CODE_VALID_FOR:           common.RequireInt64Env("AUTH_CODE_VALID_FOR"),

		DISCORD_TOKEN:  common.OptionalEnv("DISCORD_TOKEN", ""),
		SENDGRID_TOKEN: common.OptionalEnv("SENDGRID_TOKEN", ""),
	}
	return &env
}
