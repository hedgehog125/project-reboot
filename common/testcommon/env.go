package testcommon

import (
	"time"

	"github.com/hedgehog125/project-reboot/common"
)

func DefaultEnv() *common.Env {
	return &common.Env{
		IS_DEV:                        true,
		PORT:                          -1,
		MOUNT_PATH:                    "temp-test-storage",
		PROXY_ORIGINAL_IP_HEADER_NAME: "test-proxy-original-ip",

		JOB_POLL_INTERVAL:    time.Hour * 999,
		MAX_TOTAL_JOB_WEIGHT: 100,

		UNLOCK_TIME:         time.Hour * 24 * 7,
		AUTH_CODE_VALID_FOR: time.Hour * 24 * 3,

		PASSWORD_HASH_SETTINGS: &common.PasswordHashSettings{
			Time:    1,
			Memory:  1024,
			Threads: 1,
		},

		LOG_STORE_INTERVAL: time.Hour * 999,

		DISCORD_TOKEN:  "",
		SENDGRID_TOKEN: "",
	}
}
