package subfns

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/endpoints"
	"github.com/hedgehog125/project-reboot/intertypes"
)

func ConfigureServer(env *intertypes.Env) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.TrustedPlatform = env.PROXY_ORIGINAL_IP_HEADER_NAME

	engine.Static("/static", "./public")

	registerEndpoints(engine, env)

	return engine
}
func registerEndpoints(engine *gin.Engine, env *intertypes.Env) {
	endpoints.RootRedirect(engine)
}

func RunServer(engine *gin.Engine, env *intertypes.Env) {
	engine.Run(fmt.Sprintf(":%v", env.PORT))
}
