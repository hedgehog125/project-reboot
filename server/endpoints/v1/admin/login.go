package admin

import (
	"net/http"

	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/gin-gonic/gin"
)

type LoginPayload struct {
	Password string `binding:"required,min=1"         json:"password"`
	TotpCode string `binding:"required,len=6,numeric" json:"totpCode"`
}

type LoginResponse struct {
	Errors    []servercommon.ErrorDetail `binding:"required" json:"errors"`
	AdminCode string                     `                   json:"adminCode"`
}

func Login(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := LoginPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}

		adminCode, isValid := app.Core.GetAdminCode(body.Password, body.TotpCode)
		if !isValid {
			return servercommon.NewUnauthorizedError()
		}

		ginCtx.JSON(http.StatusOK, LoginResponse{
			Errors:    []servercommon.ErrorDetail{},
			AdminCode: adminCode,
		})
		return nil
	})
}
