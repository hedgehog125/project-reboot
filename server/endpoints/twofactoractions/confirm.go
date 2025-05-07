package twofactoractions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/server/servercommon"
	"github.com/hedgehog125/project-reboot/twofactoractions"
)

type ConfirmPayload struct {
	Code string `binding:"required,min=6,max=6,alphanum,lowercase" json:"code"`
}

type ConfirmResponse struct {
	Errors []string `binding:"required" json:"errors"`
}

func Confirm(app *servercommon.ServerApp) gin.HandlerFunc {
	db := app.App.Database
	clock := app.App.Clock

	return func(ctx *gin.Context) {
		body := ConfirmPayload{}
		if err := ctx.BindJSON(&body); err != nil {
			return
		}

		parsedId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ConfirmResponse{
				Errors: []string{"ID_NOT_VALID_UUID"},
			})
			return
		}

		err = twofactoractions.Confirm(parsedId, body.Code, db, clock)
		if err != nil {
			ctx.Error(servercommon.ExpectAnyOfErrors(
				err,
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
