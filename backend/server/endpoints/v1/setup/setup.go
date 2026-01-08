package setup

import (
	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
)

// Note: Some setup needs to be completed by non-setup endpoints (e.g /admin/self/messengers/...).
// For now, this is only enforced by the frontend, so once the env setup is complete,
// the setup endpoints are disabled but the main ones become enabled.
// I think this should be enough since it will be annoying enough for the admin to skip the setup that
// they'll just do it.
// And in either case, there's still a security risk if the admin gives up and leaves the setup incomplete.
func ConfigureEndpoints(group *servercommon.Group) {
	group.GET("/", GetSetup(group.App))
	if group.App.Env.ENABLE_ENV_SETUP {
		group.POST("/generate-constants", GenerateConstants(group.App))
		group.POST("/check-totp", CheckTotp(group.App))
		group.GET("/echo-headers", EchoHeaders(group.App))
	}
}
