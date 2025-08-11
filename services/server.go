package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
	"github.com/hedgehog125/project-reboot/server/endpoints"
	"github.com/hedgehog125/project-reboot/server/middleware"
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

type Server struct {
	env    *common.Env
	Router *gin.Engine
	Server *http.Server
}

func NewServer(app *common.App) *Server {
	router := gin.New()
	// router.SetTrustedProxies(nil)
	router.TrustedPlatform = app.Env.PROXY_ORIGINAL_IP_HEADER_NAME
	router.Use(gin.Logger()) // TODO: replace. This currently logs errors to std, which should probably be handled by the error middleware

	router.Static("/static", "./public")          // Has to go before otherwise files are sent but with a 404 status
	router.Use(middleware.NewTimeoutMiddleware()) // TODO: why does the error middleware have to go after? This middleware seems to write otherwise
	router.Use(middleware.NewErrorMiddleware())
	adminMiddleware := middleware.NewAdminProtectedMiddleware(app.State)
	serverApp := servercommon.ServerApp{
		App:             app,
		Router:          router,
		AdminMiddleware: adminMiddleware,
	}
	endpoints.ConfigureEndpoints(router.Group(""), &serverApp)

	return &Server{
		env:    app.Env,
		Router: router,
	}
}

func (service *Server) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", service.env.PORT),
		Handler: service.Router.Handler(),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("an error occurred while running the HTTP server:\n%v", err.Error())
		}
	}()

	service.Server = server
}
func (service *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := service.Server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the HTTP server:\n%v\n", err.Error())
	}
}
