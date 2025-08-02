package twofactoractions

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/common/dbcommon"
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
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		body := ConfirmPayload{}
		if ctxErr := servercommon.ParseBody(&body, ginCtx); ctxErr != nil {
			return ctxErr
		}
		parsedID, ctxErr := servercommon.ParseUUID(ginCtx.Param("id"))
		if ctxErr != nil {
			return ctxErr
		}

		return dbcommon.WithWriteTx(ginCtx, app.Database, func(tx *ent.Tx, ctx context.Context) error {
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
			ginCtx.JSON(http.StatusOK, ConfirmResponse{
				Errors: []servercommon.ErrorDetail{},
			})
			return nil
		})
	})
}
