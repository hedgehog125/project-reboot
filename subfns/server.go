package subfns

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/endpoints"
	"github.com/hedgehog125/project-reboot/intertypes"
)

func ConfigureServer(env *intertypes.Env) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.TrustedPlatform = env.PROXY_ORIGINAL_IP_HEADER_NAME

	engine.Use(timeout.New(timeout.WithTimeout(5 * time.Second)))

	engine.Static("/static", "./public")

	registerEndpoints(engine, env)

	return engine
}
func registerEndpoints(engine *gin.Engine, env *intertypes.Env) {
	endpoints.RootRedirect(engine)
	endpoints.RegisterUser(engine)
}

func RunServer(engine *gin.Engine, env *intertypes.Env) {
	err := engine.Run(fmt.Sprintf(":%v", env.PORT))
	if err != nil {
		log.Fatalf("error starting server. error:\n%v", err.Error())
	}
	fmt.Printf("running on port %v", env.PORT)
}
