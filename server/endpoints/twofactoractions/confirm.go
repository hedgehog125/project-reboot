package twofactoractions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type ConfirmPayload struct {
	Code string `binding:"required,min=9,max=9,alphanum,lowercase" json:"code"`
}

type ConfirmResponse struct {
	Errors []servercommon.ErrorDetail `binding:"required" json:"errors"`
}

func Confirm(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.WithTx(app, func(ctx *gin.Context, tx *ent.Tx) error {
		body := ConfirmPayload{}
		if ctxErr := servercommon.ParseBody(&body, ctx); ctxErr != nil {
			return ctxErr
		}
		parsedID, ctxErr := servercommon.ParseUUID(ctx.Param("id"))
		if ctxErr != nil {
			return ctxErr
		}

		_, commErr := app.TwoFactorActions.Confirm(parsedID, body.Code, ctx)
		if commErr != nil {
			return servercommon.ExpectAnyOfErrors(
				commErr,
				[]error{
					twofactoractions.ErrNotFound,
					twofactoractions.ErrExpired,
					twofactoractions.ErrWrongCode,
				},
				http.StatusUnauthorized, nil,
			)
		}

		ctx.JSON(http.StatusOK, ConfirmResponse{
			Errors: []servercommon.ErrorDetail{},
		})
		return nil
	})
}
