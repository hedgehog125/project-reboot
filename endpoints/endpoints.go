package endpoints

import (
	"github.com/gin-gonic/gin"
)

func RootRedirect(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		ctx.File("./public/index.html")
	})
}

type RegisterUserPayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
	Content  string `json:"content"  binding:"required,min=1,max=100000000"` // 100 MB but base64 encoded
	MimeType string `json:"mimeType" binding:"required,min=1,max=256"`
}

func RegisterUser(engine *gin.Engine) {
	engine.POST("/api/v1/register", func(ctx *gin.Context) {
		body := RegisterUserPayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}
	})
}
