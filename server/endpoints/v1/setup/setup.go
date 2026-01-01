package setup

import (
	"github.com/NicoClack/cryptic-stash/server/servercommon"
)

func ConfigureEndpoints(group *servercommon.Group) {
	group.POST("/generate-constants", GenerateConstants(group.App))
	group.POST("/check-totp", CheckTotp(group.App))
	group.GET("/echo-headers", EchoHeaders(group.App))
}
