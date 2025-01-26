package subfns

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/endpoints"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/intertypes"
)

func ConfigureServer(state *intertypes.State, dbClient *ent.Client, env *intertypes.Env) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.TrustedPlatform = env.PROXY_ORIGINAL_IP_HEADER_NAME

	engine.Use(endpoints.NewTimeoutMiddleware())
	adminMiddleware := endpoints.NewAdminProtectedMiddleware(state)

	engine.Static("/static", "./public")

	registerEndpoints(engine, adminMiddleware, dbClient, env)

	return engine
}
func registerEndpoints(engine *gin.Engine, adminMiddleware gin.HandlerFunc, dbClient *ent.Client, env *intertypes.Env) {
	endpoints.RootRedirect(engine)
	endpoints.RegisterUser(engine, adminMiddleware, dbClient)
}

func RunServer(engine *gin.Engine, env *intertypes.Env) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", env.PORT),
		Handler: engine.Handler(),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("an error occurred while running the HTTP server:\n%v", err.Error())
		}
	}()

	return server
}
func ShutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the HTTP server:\n%v\n", err.Error())
	}
}
