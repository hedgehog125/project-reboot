package servercommon

import (
	"github.com/gin-gonic/gin"
)

func NewHandler(
	handler func(ginCtx *gin.Context) error,
) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		stdErr := handler(ginCtx)
		if stdErr != nil {
			ginCtx.Error(stdErr)
		}
	}
}
