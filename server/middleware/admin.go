package middleware

import (
	"net/http"
	"strings"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

func NewAdminProtected(core common.CoreService) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		headerValue := ginCtx.GetHeader("authorization")
		if headerValue == "" {
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "authorization header is required",
						Code:    "MISSING_AUTHORIZATION_HEADER",
					},
				},
			})
			return
		}
		headerParts := strings.SplitN(headerValue, " ", 2)
		if len(headerParts) != 2 {
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "malformed authorization header",
						Code:    "MALFORMED_AUTHORIZATION_HEADER",
					},
				},
			})
			return
		}
		if headerParts[0] != "AdminCode" {
			ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "unsupported authorization scheme",
						Code:    "UNSUPPORTED_AUTHORIZATION_SCHEME",
					},
				},
			})
			return
		}

		if !core.CheckAdminCode(headerParts[1]) {
			ginCtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []servercommon.ErrorDetail{
					{
						Message: "invalid admin code",
						Code:    "INVALID_ADMIN_CODE",
					},
				},
			})
			return
		}
		ginCtx.Next()
	}
}
