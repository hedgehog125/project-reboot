package servercommon

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseBody(obj any, ctx *gin.Context) *Error {
	// TODO: this is not secure. The error type should is a validator.ValidationErrors which should be processed
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return NewError(ErrWrapperParseBodyJson.Wrap(err)).
			SetStatus(http.StatusBadRequest).
			AddDetail(ErrorDetail{
				Message: err.Error(),
				Code:    "INVALID_BODY_JSON",
			}).
			DisableLogging()
	}
	return nil
}
