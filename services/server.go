package services

import (
	"context"
	"errors"
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
	App    *common.App
	Router *gin.Engine
	Server *http.Server
}

func NewServer(app *common.App) *Server {
	router := gin.New()
	if app.Env.PROXY_ORIGINAL_IP_HEADER_NAME == "" {
		router.SetTrustedProxies(nil)
	} else {
		router.TrustedPlatform = app.Env.PROXY_ORIGINAL_IP_HEADER_NAME
	}
	router.Use(middleware.NewLogger(app.Logger))
	router.Use(gin.Logger()) // TODO: make the custom logger log completed requests so this isn't needed
	// TODO: ^ this is logging "Error #01: ..."

	router.LoadHTMLGlob("./server/templates/*.html")
	router.Use(middleware.NewRateLimiting("api", app.RateLimiter))
	router.Static("/static", "./public") // Has to go before otherwise files are sent but with a 404 status
	router.Use(middleware.NewTimeout())  // TODO: why does the error middleware have to go after? This middleware seems to write otherwise
	router.Use(middleware.NewError())
	adminMiddleware := middleware.NewAdminProtected(app.Core)
	serverApp := &servercommon.ServerApp{
		App:             app,
		Router:          router,
		AdminMiddleware: adminMiddleware,
	}
	endpoints.ConfigureEndpoints(router.Group(""), serverApp) // TODO: rework to be more like jobs registry, embed *gin.RouterGroup

	return &Server{
		App:    app,
		Router: router,
	}
}

func (service *Server) Start() {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", service.App.Env.PORT),
		Handler:           service.Router.Handler(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	service.Server = server

	go func() {
		stdErr := server.ListenAndServe()
		if stdErr != nil && !errors.Is(stdErr, http.ErrServerClosed) {
			log.Fatalf("an error occurred while starting the HTTP server:\n%v", stdErr.Error())
		}
	}()
}
func (service *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stdErr := service.Server.Shutdown(ctx)
	if stdErr != nil {
		service.App.Logger.Warn("an error occurred while shutting down the HTTP server", stdErr)
	}
}
