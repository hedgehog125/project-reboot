package services

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/server"
	"github.com/NicoClack/cryptic-stash/server/endpoints"
	"github.com/NicoClack/cryptic-stash/server/middleware"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"errors": []string{},
		})
	})
	router.Use(middleware.NewTimeout())

	router.LoadHTMLFS(http.FS(server.TemplateFiles.FS), fmt.Sprintf("%v/*.html", server.TemplateFiles.Path))
	router.Use(middleware.NewRateLimiting("api", app.RateLimiter))

	embeddedFolder, stdErr := static.EmbedFolder(server.PublicFiles.FS, server.PublicFiles.Path)
	if stdErr != nil {
		log.Fatalf("failed to load embedded public files:\n%v", stdErr.Error())
	}
	router.StaticFS("/", embeddedFolder)

	router.Use(middleware.NewError())

	adminMiddleware := middleware.NewAdminProtected(app.Core)
	serverApp := &servercommon.ServerApp{
		App:             app,
		Router:          router,
		AdminMiddleware: adminMiddleware,
	}
	endpoints.ConfigureEndpoints(&servercommon.Group{
		RouterGroup: router.Group(""),
		App:         serverApp,
	})

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%v", app.Env.PORT),
		Handler:           router.Handler(),
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return &Server{
		App:    app,
		Router: router,
		Server: httpServer,
	}
}

func (service *Server) Start() {
	go func() {
		stdErr := service.Server.ListenAndServe()
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

func (service *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service.Server.Handler.ServeHTTP(w, r)
}
