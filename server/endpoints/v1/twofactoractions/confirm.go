package twofactoractions

import (
	"context"
	"net/http"

	"github.com/NicoClack/cryptic-stash/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/ent"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/NicoClack/cryptic-stash/twofactoractions"
	"github.com/gin-gonic/gin"
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
			_, wrappedErr := app.TwoFactorActions.Confirm(parsedID, body.Code, ctx)
			if wrappedErr != nil {
				return servercommon.ExpectAnyOfErrors(
					wrappedErr,
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
