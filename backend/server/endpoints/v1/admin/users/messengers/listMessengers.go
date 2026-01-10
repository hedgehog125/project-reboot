package messengers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/NicoClack/cryptic-stash/backend/ent/user"
	"github.com/NicoClack/cryptic-stash/backend/server/servercommon"
	"github.com/gin-gonic/gin"
)

type ListMessengersResponse struct {
	Errors     []servercommon.ErrorDetail `binding:"required" json:"errors"`
	Messengers map[string]*Messenger      `binding:"required" json:"messengers"`
}

type Messenger struct {
	Name          string          `binding:"required" json:"name"`
	Created       bool            `binding:"required" json:"created"`
	Enabled       bool            `binding:"required" json:"enabled"`
	Options       json.RawMessage `binding:"required" json:"options"`
	OptionsSchema json.RawMessage `binding:"required" json:"optionsSchema"`
}

func ListMessengers(app *servercommon.ServerApp) gin.HandlerFunc {
	return servercommon.NewHandler(func(ginCtx *gin.Context) error {
		userID, ctxErr := servercommon.ParseObjectID(ginCtx.Param("id"))
		if ctxErr != nil {
			return ctxErr
		}
		userOb, stdErr := dbcommon.WithReadTx(
			ginCtx.Request.Context(), app.Database,
			func(tx *ent.Tx, ctx context.Context) (*ent.User, error) {
				return tx.User.Query().
					Where(user.ID(userID)).
					WithMessengers().
					Only(ctx)
			},
		)
		if stdErr != nil {
			return servercommon.Send404IfNotFound(stdErr)
		}

		definitions := app.Messengers.AllPublicDefinitions()
		userMessengers := make(map[string]*Messenger, len(definitions))
		for _, messengerOb := range userOb.Edges.Messengers {
			versionedType := common.GetVersionedType(messengerOb.Type, messengerOb.Version)
			definition, ok := app.Messengers.GetPublicDefinition(versionedType)
			if !ok {
				return fmt.Errorf(
					"user %v has %v messenger configured but it has no definition",
					userOb.ID,
					versionedType,
				)
			}

			//exhaustruct:enforce
			userMessengers[versionedType] = &Messenger{
				Name:          definition.Name,
				Created:       true,
				Enabled:       messengerOb.Enabled,
				Options:       messengerOb.Options,
				OptionsSchema: definition.OptionsSchema,
			}
		}
		for _, definition := range definitions {
			versionedType := common.GetVersionedType(definition.ID, definition.Version)
			_, ok := userMessengers[versionedType]
			if ok {
				continue
			}

			//exhaustruct:enforce
			userMessengers[versionedType] = &Messenger{
				Name:          definition.Name,
				Created:       false,
				Enabled:       false,
				Options:       nil,
				OptionsSchema: definition.OptionsSchema,
			}
		}

		//exhaustruct:enforce
		ginCtx.JSON(http.StatusOK, ListMessengersResponse{
			Errors:     []servercommon.ErrorDetail{},
			Messengers: userMessengers,
		})
		return nil
	})
}
