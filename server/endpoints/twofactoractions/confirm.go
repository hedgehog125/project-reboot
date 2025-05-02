package twofactoractions

import (
	"context"
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
	dbClient := app.App.Database.Client()
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

		action, err := dbClient.TwoFactorAction.Get(context.Background(), parsedId)
		if err != nil {
			ctx.Error(servercommon.Send404IfNotFound(err))
			return
		}

		if action.ExpiresAt.Before(clock.Now()) {
			err := dbClient.TwoFactorAction.DeleteOne(action).Exec(ctx)
			if err != nil {
				// TODO: log warning
			}
			ctx.JSON(http.StatusNotFound, ConfirmResponse{
				Errors: []string{},
			})
			return
		}

		// TODO: validate body.Code!
		// TODO: execute the action
		err = twofactoractions.Confirm(action)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, ConfirmResponse{
			Errors: []string{},
		})
	}
}
