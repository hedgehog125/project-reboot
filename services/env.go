package services

import (
	"log"
	"log/slog"
	"os"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() *common.Env {
	_, isDevEnvDefined := os.LookupEnv("IS_DEV")
	if !isDevEnvDefined {
		stdErr := godotenv.Load(".env")
		if stdErr != nil {
			// The usual logger hasn't been created yet
			slog.Warn("error loading .env file", "error", stdErr)
		}
	}

	env := &common.Env{
		IS_DEV:                        common.RequireBoolEnv("IS_DEV"),
		PORT:                          common.RequireIntEnv("PORT"),
		MOUNT_PATH:                    common.RequireEnv("MOUNT_PATH"),
		PROXY_ORIGINAL_IP_HEADER_NAME: common.RequireEnv("PROXY_ORIGINAL_IP_HEADER_NAME"),

		JOB_POLL_INTERVAL:    common.RequireSecondsEnv("JOB_POLL_INTERVAL"),
		MAX_TOTAL_JOB_WEIGHT: common.RequireIntEnv("MAX_TOTAL_JOB_WEIGHT"),

		UNLOCK_TIME:              common.RequireSecondsEnv("UNLOCK_TIME"),
		AUTH_CODE_VALID_FOR:      common.RequireSecondsEnv("AUTH_CODE_VALID_FOR"),
		USED_AUTH_CODE_VALID_FOR: common.RequireSecondsEnv("USED_AUTH_CODE_VALID_FOR"),
		PASSWORD_HASH_SETTINGS: &common.PasswordHashSettings{
			Time:    common.RequireUint32Env("PASSWORD_HASH_TIME"),
			Memory:  common.RequireUint32Env("PASSWORD_HASH_MEMORY"),
			Threads: common.RequireUint8Env("PASSWORD_HASH_THREADS"),
		},

		LOG_STORE_INTERVAL:    common.RequireMillisecondsEnv("LOG_STORE_INTERVAL"),
		ADMIN_USERNAME:        common.RequireEnv("ADMIN_USERNAME"),
		ADMIN_MESSAGE_TIMEOUT: common.RequireSecondsEnv("ADMIN_MESSAGE_TIMEOUT"),

		DISCORD_TOKEN:  common.OptionalEnv("DISCORD_TOKEN", ""),
		SENDGRID_TOKEN: common.OptionalEnv("SENDGRID_TOKEN", ""),
	}
	ValidateEnvironmentVariables(env)
	return env
}
func ValidateEnvironmentVariables(env *common.Env) {
	if float64(env.AUTH_CODE_VALID_FOR)/float64(env.UNLOCK_TIME) < 1.1 {
		log.Fatalf(
			"AUTH_CODE_VALID_FOR must be at least slightly larger than UNLOCK_TIME because a download requires the auth code to be valid and the unlock time needs to have passed",
		)
	}
}
