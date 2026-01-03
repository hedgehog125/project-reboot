package middleware

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
	"github.com/gin-gonic/gin"
)

func NewStaticFS(embeddedFS embed.FS, prefix string) gin.HandlerFunc {
	subFS, stdErr := fs.Sub(embeddedFS, prefix)
	if stdErr != nil {
		log.Fatalf("error calling fs.Sub:\n%v", stdErr.Error())
	}
	fileServer := http.FileServerFS(subFS)

	return func(ginCtx *gin.Context) {
		if strings.HasPrefix(ginCtx.Request.URL.Path, "/api") {
			ginCtx.JSON(http.StatusNotFound, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Code:    "ENDPOINT_NOT_FOUND",
						Message: "Endpoint not found",
					},
				},
			})
			return
		}

		// We can make things slightly more efficient and this 200 fallback behaviour slightly less confusing by only doing
		// local redirects for pages, not resources
		// Also this means fileServer can redirect /index.html to /
		if path.Ext(ginCtx.Request.URL.Path) == "" {
			filePath := strings.TrimPrefix(path.Clean(ginCtx.Request.URL.Path), "/")
			if filePath == "" {
				filePath = "."
			}
			file, stdErr := subFS.Open(filePath)
			if stdErr == nil {
				_ = file.Close()
			} else {
				ginCtx.Request.URL.Path = "/"
			}
		}

		fileServer.ServeHTTP(ginCtx.Writer, ginCtx.Request)
		ginCtx.Abort()
	}
}
