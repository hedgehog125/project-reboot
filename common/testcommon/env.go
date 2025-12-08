package testcommon

import (
	"time"

	"github.com/NicoClack/cryptic-stash/common"
)

func DefaultEnv() *common.Env {
	//exhaustruct:enforce
	return &common.Env{
		IS_DEV:                        true,
		PORT:                          -1,
		MOUNT_PATH:                    "temp-test-storage",
		PROXY_ORIGINAL_IP_HEADER_NAME: "test-proxy-original-ip",
		CLEAN_UP_INTERVAL:             time.Hour,
		FULL_GC_INTERVAL:              0,

		JOB_POLL_INTERVAL:    time.Hour * 999,
		MAX_TOTAL_JOB_WEIGHT: 100,

		UNLOCK_TIME:                      time.Hour * 24 * 7,
		AUTH_CODE_VALID_FOR:              time.Hour * 24 * 3,
		USED_AUTH_CODE_VALID_FOR:         time.Hour,
		ACTIVE_SESSION_REMINDER_INTERVAL: time.Hour * 24,
		MIN_SUCCESSFUL_MESSAGE_COUNT:     1,

		PASSWORD_HASH_SETTINGS: &common.PasswordHashSettings{
			Time:    1,
			Memory:  1024,
			Threads: 1,
		},

		LOG_STORE_INTERVAL:    time.Hour * 999,
		ADMIN_USERNAME:        "", // The test will need to set this up if required
		ADMIN_MESSAGE_TIMEOUT: time.Minute,
		MIN_ADMIN_MESSAGE_GAP: time.Minute * 5,
		MIN_CRASH_SIGNAL_GAP:  time.Hour * 24,
		PANIC_ON_ERROR:        true,

		ENABLE_DEVELOP_MESSENGER: false,
		DISCORD_TOKEN:            "",
		SENDGRID_TOKEN:           "",
	}
}
