package middleware

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func NewStaticFS(embeddedFS embed.FS, prefix string) gin.HandlerFunc {
	subFS, stdErr := fs.Sub(embeddedFS, prefix)
	if stdErr != nil {
		log.Fatalf("error calling fs.Sub:\n%v", stdErr.Error())
	}
	fileServer := http.FileServerFS(subFS)

	return func(ginCtx *gin.Context) {
		// We can make things slightly more efficient and this 200 fallback behaviour slightly less confusing by only doing
		// local redirects for pages, not resources
		// Also this means fileServer can redirect /index.html to /
		if path.Ext(ginCtx.Request.URL.Path) == "" {
			file, stdErr := subFS.Open(path.Clean(ginCtx.Request.URL.Path))
			if stdErr == nil {
				file.Close()
			} else {
				ginCtx.Request.URL.Path = "/"
			}
		}

		fileServer.ServeHTTP(ginCtx.Writer, ginCtx.Request)
		ginCtx.Abort()
	}
}
