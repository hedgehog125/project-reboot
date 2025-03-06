package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedgehog125/project-reboot/core"
	"github.com/hedgehog125/project-reboot/ent"
	"github.com/hedgehog125/project-reboot/ent/user"
	"github.com/hedgehog125/project-reboot/intertypes"
	"github.com/hedgehog125/project-reboot/util"
	"github.com/jonboulle/clockwork"
)

type GetAuthorizationCodePayload struct {
	Username string `json:"username" binding:"required,min=1,max=32,alphanum,lowercase"`
	Password string `json:"password" binding:"required,min=8,max=256"`
}

type GetAuthorizationCodeResponse struct {
	Errors                   []string  `json:"errors" binding:"required"`
	AuthorizationCode        string    `json:"authorizationCode"`
	AuthorizationCodeValidAt time.Time `json:"authorizationCodeValidAt"`
}

// TODO: split into 2 different endpoints
func GetAuthorizationCode(engine *gin.Engine, dbClient *ent.Client, clock clockwork.Clock, env *intertypes.Env) gin.HandlerFunc {
	sendUnauthorizedError := func(ctx *gin.Context) {
		ctx.JSON(http.StatusUnauthorized, GetAuthorizationCodeResponse{
			Errors: []string{"INCORRECT_USERNAME_OR_PASSWORD_OR_AUTH_CODE"},
		})
	}

	return func(ctx *gin.Context) {
		body := GetAuthorizationCodePayload{}
		if err := ctx.BindJSON(&body); err != nil { // TODO: request size limits?
			return
		}

		userRow, err := dbClient.User.Query().
			Where(user.Username(body.Username)).
			Select(user.FieldPasswordHash, user.FieldPasswordSalt, user.FieldHashTime, user.FieldHashMemory, user.FieldHashKeyLen).
			Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				sendUnauthorizedError(ctx)
			} else {
				fmt.Printf("warning: an error occurred while reading user data:\n%v\n", err.Error())
				ctx.JSON(http.StatusInternalServerError, GetAuthorizationCodeResponse{
					Errors: []string{"INTERNAL"},
				})
			}

			return
		}

		authorized := core.CheckPassword(
			body.Password,
			userRow.PasswordHash,
			userRow.PasswordSalt,
			&core.HashSettings{
				Time:   userRow.HashTime,
				Memory: userRow.HashMemory,
				KeyLen: userRow.HashKeyLen,
			},
		)

		var authCode []byte = nil
		if authorized {
			authCode = util.CryptoRandomBytes(core.AUTH_CODE_BYTE_LENGTH)
		}
		validAt := clock.Now().UTC().
			Add(time.Duration(env.UNLOCK_TIME) * time.Second)

		_, err = dbClient.LoginAttempt.Create().
			SetUser(userRow).
			SetCode(authCode).
			SetCodeValidFrom(validAt).
			SetInfo(&intertypes.LoginAttemptInfo{
				UserAgent: ctx.Request.UserAgent(),
				IP:        ctx.ClientIP(),
			}).Save(context.Background())
		if err != nil {
			fmt.Printf("warning: an error occurred while creating a login attempt:\n%v\n", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"errors": []string{"INTERNAL"},
			})
			return
		}

		if !authorized {
			sendUnauthorizedError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, GetAuthorizationCodeResponse{
			Errors:                   []string{},
			AuthorizationCode:        base64.StdEncoding.EncodeToString(authCode),
			AuthorizationCodeValidAt: validAt,
		})
	}
}
