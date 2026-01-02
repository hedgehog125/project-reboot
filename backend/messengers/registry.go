package messengers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/NicoClack/cryptic-stash/backend/common/dbcommon"
	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/NicoClack/cryptic-stash/backend/jobs"
)

const JobNamePrefix = "messengers"

type Registry struct {
	App        *common.App
	messengers map[string]*Definition
	// RegisterJobs can be called with a prefix, so this is JobNamePrefix + that
	jobNamePrefix string
}

type Definition struct {
	ID      string
	Version int
	// Supplemental messengers often can't tell if a message was successfully sent. So they increase the
	// chance the user is notified but are ignored when assessing if the user was sufficiently notified.
	// This also means that the user needs to have at least one non-supplemental messenger configured,
	// otherwise they'll never be considered sufficiently notified.
	IsSupplemental bool
	// Returns the data the Handler needs, typically a struct containing the formatted message and
	// some sort of contact (e.g a username)
	// If the user doesn't have the right contacts for this messenger, return messengers.ErrNoContactForUser.Clone()
	Prepare PrepareFunc
	// The return type of Prepare
	BodyType          any
	Handler           HandlerFunc
	jobDefinition     *jobs.Definition
	reflectedBodyType reflect.Type
}

type PrepareFunc = func(message *common.Message) (any, error)
type HandlerFunc func(messengerCtx *Context) error

type JobContext = jobs.Context
type Context struct {
	*JobContext
	confirmedSent bool
}

func (ctx *Context) ConfirmSent() {
	ctx.confirmedSent = true
}

func NewRegistry(app *common.App) *Registry {
	return &Registry{
		messengers: make(map[string]*Definition),
		App:        app,
	}
}

type bodyWrapperType struct {
	MessageType            common.MessageType
	VersionedMessengerType string
	SessionIDs             []int
	Inner                  string
}

func (registry *Registry) Register(definition *Definition) {
	versionedType := common.GetVersionedType(definition.ID, definition.Version)
	if _, exists := registry.messengers[versionedType]; exists {
		log.Fatalf("messenger definition with ID \"%s\" already exists", versionedType)
	}
	definition.reflectedBodyType = reflect.TypeOf(definition.BodyType)
	jobs.AssertTypeIsValidBodyType(definition.reflectedBodyType, versionedType)

	definition.jobDefinition = &jobs.Definition{
		ID:      definition.ID,
		Version: definition.Version,
		Handler: func(jobCtx *jobs.Context) error {
			body := &bodyWrapperType{}
			wrappedErr := jobCtx.Decode(body)
			if wrappedErr != nil {
				return ErrWrapperHandlerWrapper.Wrap(wrappedErr)
			}

			newJobCtx := *jobCtx
			newJobCtx.Body = json.RawMessage(body.Inner)
			messengerCtx := &Context{
				JobContext:    &newJobCtx,
				confirmedSent: false,
			}
			stdErr := definition.Handler(messengerCtx)
			if stdErr != nil {
				if body.MessageType == common.MessageAdminError {
					wrappedErr := common.WrapErrorWithCategories(stdErr)
					if wrappedErr.MaxRetries() > 0 || wrappedErr.MaxRetries() == -1 {
						if registry.App.Clock.Since(jobCtx.OriginallyDueAt) >= registry.App.Env.ADMIN_MESSAGE_TIMEOUT {
							jobCtx.Logger.ErrorContext(
								context.WithValue(context.Background(), common.AdminNotificationFallbackKey{}, true),
								"failed to notify admin about an error before ADMIN_MESSAGE_TIMEOUT, "+
									"will now possibly crash to notify them earlier",
								"jobID",
								jobCtx.ID,
							)
						}
					} else {
						jobCtx.Logger.ErrorContext(
							context.WithValue(context.Background(), common.AdminNotificationFallbackKey{}, true),
							"failed to notify admin about an error! will now possibly crash to notify them earlier",
							"jobID",
							jobCtx.ID,
						)
					}
				}
				return stdErr
			}

			if body.MessageType == common.MessageLogin ||
				body.MessageType == common.MessageActiveSessionReminder {
				// This is in a separate transaction, so we could successfully commit the messenger's transaction
				// but roll back this one. But that's ok since it's best to undercount the successful login alerts sent
				stdErr := dbcommon.WithWriteTx(
					jobCtx.Context, registry.App.Database,
					func(tx *ent.Tx, ctx context.Context) error {
						return tx.LoginAlert.MapCreateBulk(
							body.SessionIDs,
							func(alertCreate *ent.LoginAlertCreate, i int) {
								alertCreate.
									SetSentAt(registry.App.Clock.Now()).
									SetVersionedMessengerType(body.VersionedMessengerType).
									SetConfirmed(messengerCtx.confirmedSent).
									SetSessionID(body.SessionIDs[i])
							},
						).Exec(ctx)
					},
				)
				if stdErr != nil {
					// TODO: handle missing session IDs, stop this being atomic
					jobCtx.Logger.Error(
						"failed to create LoginAlert objects for successfully sent message, if not enough objects are created, "+
							"the user won't be able to download their data once their session becomes valid",
						"error",
						stdErr,
						"sessionIDs",
						body.SessionIDs,
					)
				}
			}
			return nil
		},
		BodyType: &bodyWrapperType{},
		Weight:   1,
	}
	registry.messengers[versionedType] = definition
}
func (registry *Registry) RegisterJobs(group *jobs.RegistryGroup) {
	registry.jobNamePrefix = common.JoinPaths(group.Path, JobNamePrefix)
	prefixedGroup := group.Group(JobNamePrefix)
	for _, messenger := range registry.messengers {
		prefixedGroup.Register(messenger.jobDefinition)
	}
}

func (registry *Registry) Send(
	versionedType string, message *common.Message,
	sendTime time.Time,
	ctx context.Context,
) common.WrappedError {
	messengerDef, ok := registry.messengers[versionedType]
	if !ok {
		return ErrWrapperSend.Wrap(ErrUnknownMessengerType)
	}
	tx := ent.TxFromContext(ctx)
	if tx == nil {
		return ErrWrapperSend.Wrap(common.ErrNoTxInContext)
	}

	preparedData, stdErr := messengerDef.Prepare(message)
	if stdErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperPrepare.Wrap(stdErr))
	}
	bodyType := reflect.TypeOf(preparedData)
	if bodyType != messengerDef.reflectedBodyType {
		return ErrWrapperSend.Wrap(ErrWrapperEnqueueJob.Wrap(
			jobs.ErrWrapperEncode.Wrap(jobs.ErrWrapperInvalidBody.Wrap(
				fmt.Errorf("body type %s isn't the expected type %s",
					preparedData, messengerDef.reflectedBodyType),
			)),
		))
	}
	encoded, stdErr := json.Marshal(preparedData)
	if stdErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperEnqueueJob.Wrap(
			jobs.ErrWrapperEncode.Wrap(jobs.ErrWrapperInvalidBody.Wrap(stdErr)),
		))
	}

	_, wrappedErr := registry.App.Jobs.EnqueueWithModifier(
		common.JoinPaths(registry.jobNamePrefix, versionedType),
		&bodyWrapperType{
			MessageType:            message.Type,
			VersionedMessengerType: versionedType,
			SessionIDs:             message.SessionIDs,
			Inner:                  string(encoded),
		},
		func(jobCreate *ent.JobCreate) {
			jobCreate.SetDueAt(sendTime)
		},
		ctx,
	)
	if wrappedErr != nil {
		return ErrWrapperSend.Wrap(ErrWrapperEnqueueJob.Wrap(wrappedErr))
	}

	logMessage, wrappedErr := FormatDefaultMessage(message)
	logger := registry.App.Logger.With(
		"userID", message.User.ID,
		"messageType", versionedType,
	)
	if wrappedErr == nil {
		logger.Info(
			"sending message to user",
			"message", logMessage,
		)
	} else {
		logger.Warn(
			"sending a message that FormatDefaultMessage couldn't format",
			"error", wrappedErr,
		)
	}

	return nil
}
func (registry *Registry) SendUsingAll(
	message *common.Message,
	sendTime time.Time,
	ctx context.Context,
) (int, map[string]common.WrappedError, common.WrappedError) {
	errs := make(map[string]common.WrappedError)
	messagesQueued := 0
	for versionedType := range registry.messengers {
		wrappedErr := registry.Send(versionedType, message, sendTime, ctx)
		if wrappedErr == nil {
			messagesQueued++
		} else {
			// TODO: remove this ErrNoContactForUser sentinel error and have a separate hook
			// to check if the messenger has the right contacts?
			errs[versionedType] = wrappedErr
			if !ErrWrapperPrepare.HasWrapped(wrappedErr) {
				return messagesQueued, errs, ErrWrapperSendUsingAll.Wrap(wrappedErr)
			}
			if !errors.Is(wrappedErr, ErrNoContactForUser) {
				// Just log an error and let the admin deal with this, there's not much the user can do
				registry.App.Logger.Error(
					"failed to prepare message",
					"messengerType", versionedType,
					"error", wrappedErr,
				)
			}
		}
	}
	return messagesQueued, errs, nil
}

// Note: lastSendTime will be zero for the first call
func (registry *Registry) SendBulk(
	messages []*common.Message, sendTimeFunc func(lastSendTime time.Time, index int) time.Time, ctx context.Context,
) common.WrappedError {
	if len(messages) == 0 {
		return nil
	}
	sendTime := sendTimeFunc(time.Time{}, 0)
	for index, message := range messages {
		_, _, wrappedErr := registry.SendUsingAll(
			message, sendTime, ctx,
		)
		if wrappedErr != nil {
			return ErrWrapperSendBulk.Wrap(wrappedErr)
		}
		sendTime = sendTimeFunc(sendTime, index+1)
	}
	return nil
}

func (registry *Registry) GetConfiguredMessengerTypes(user *ent.User) []string {
	configuredTypes := []string{}
	for versionedType, messengerDef := range registry.messengers {
		_, wrappedErr := messengerDef.Prepare(&common.Message{
			Type: common.MessageTest,
			User: user,
		})
		if wrappedErr != nil {
			continue
		}
		configuredTypes = append(configuredTypes, versionedType)
	}
	return configuredTypes
}
func (registry *Registry) GetPublicDefinition(versionedType string) (*common.MessengerDefinition, bool) {
	messengerDef, ok := registry.messengers[versionedType]
	if !ok {
		return nil, false
	}
	return &common.MessengerDefinition{
		ID:             messengerDef.ID,
		Version:        messengerDef.Version,
		IsSupplemental: messengerDef.IsSupplemental,
	}, true
}
