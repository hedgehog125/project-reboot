package setup

import (
	"net/http"

	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

type EchoHeadersResponse struct {
	Errors  []servercommon.ErrorDetail `binding:"required" json:"errors"`
	Headers map[string][]string        `                   json:"headers"`
}

func EchoHeaders(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		ginCtx.JSON(http.StatusOK, EchoHeadersResponse{
			Errors:  []servercommon.ErrorDetail{},
			Headers: ginCtx.Request.Header,
		})
		return nil
	})
}
