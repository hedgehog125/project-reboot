package endpoints

import (
	"github.com/gin-gonic/gin"
)

func RootRedirect(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	})
}
