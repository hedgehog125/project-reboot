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
	"github.com/hedgehog125/project-reboot/server/servercommon"
)

func NewServer(app *common.App) common.ServerService {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.TrustedPlatform = app.Env.PROXY_ORIGINAL_IP_HEADER_NAME

	router.Static("/static", "./public") // Has to go before otherwise files are sent but with a 404 status
	router.Use(endpoints.NewTimeoutMiddleware())
	adminMiddleware := endpoints.NewAdminProtectedMiddleware(app.State)
	serverApp := servercommon.ServerApp{
		App:             app,
		Router:          router,
		AdminMiddleware: adminMiddleware,
	}
	endpoints.ConfigureEndpoints(router.Group(""), &serverApp)

	return &serverService{
		env:    app.Env,
		router: router,
	}
}

type serverService struct {
	env    *common.Env
	router *gin.Engine
	server *http.Server
}

func (service *serverService) Start() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", service.env.PORT),
		Handler: service.router.Handler(),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("an error occurred while running the HTTP server:\n%v", err.Error())
		}
	}()

	service.server = server
}
func (service *serverService) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := service.server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("warning: an error occurred while shutting down the HTTP server:\n%v\n", err.Error())
	}
}
