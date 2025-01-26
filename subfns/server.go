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
	"github.com/jonboulle/clockwork"
)

func ConfigureServer(
	state *intertypes.State,
	dbClient *ent.Client,
	clock clockwork.Clock,
	env *intertypes.Env,
) *gin.Engine {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.TrustedPlatform = env.PROXY_ORIGINAL_IP_HEADER_NAME

	engine.Static("/static", "./public") // Has to go before otherwise files are sent but with a 404 status
	engine.Use(endpoints.NewTimeoutMiddleware())
	adminMiddleware := endpoints.NewAdminProtectedMiddleware(state)

	registerEndpoints(engine, adminMiddleware, dbClient, clock, env)

	return engine
}
func registerEndpoints(
	engine *gin.Engine,
	adminMiddleware gin.HandlerFunc,
	dbClient *ent.Client,
	clock clockwork.Clock,
	env *intertypes.Env,
) {
	endpoints.RootRedirect(engine)
	endpoints.RegisterUser(engine, adminMiddleware, dbClient)
	endpoints.GetUserDownload(engine, dbClient, clock, env)
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
