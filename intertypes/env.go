package intertypes

type Env struct {
	PORT                          int
	MOUNT_PATH                    string
	PROXY_ORIGINAL_IP_HEADER_NAME string
	UNLOCK_TIME                   int64 // In seconds
	// TODO: implement
	AUTH_CODE_VALID_FOR int64 // In seconds
}
