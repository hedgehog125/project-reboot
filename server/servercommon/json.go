package servercommon

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common"
)

func ParseBody(obj any, ctx *gin.Context) *ContextError {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return &ContextError{
			Err:        err,
			Status:     http.StatusBadRequest,
			ErrorCodes: []string{"INVALID_BODY"},
			Category:   common.ErrTypeClient,
			ShouldLog:  false,
		}
	}
	return nil
}
