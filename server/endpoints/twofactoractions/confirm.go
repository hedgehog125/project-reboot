package twofactoractions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type ConfirmPayload struct {
	Code string `binding:"required,min=9,max=9,alphanum,lowercase" json:"code"`
}

type ConfirmResponse struct {
	Errors []string `binding:"required" json:"errors"`
}

func Confirm(app *servercommon.ServerApp) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := ConfirmPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			ctx.Error(ctxErr)
			return
		}

		parsedId, stdErr := uuid.Parse(ctx.Param("id"))
		if stdErr != nil {
			ctx.JSON(http.StatusBadRequest, ConfirmResponse{
				Errors: []string{"ID_NOT_VALID_UUID"},
			})
			return
		}

		commErr := app.App.TwoFactorAction.Confirm(parsedId, body.Code)
		if commErr != nil {
			ctx.Error(servercommon.ExpectAnyOfErrors(
				commErr,
				[]error{
					twofactoractions.ErrNotFound,
					twofactoractions.ErrExpired,
					twofactoractions.ErrWrongCode,
				},
				http.StatusUnauthorized, "",
			))
			return
		}

		ctx.JSON(http.StatusOK, ConfirmResponse{
			Errors: []string{},
		})
	}
}
