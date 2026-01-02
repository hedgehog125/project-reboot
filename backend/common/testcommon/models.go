package testcommon

import (
	"context"
	"fmt"

	"github.com/NicoClack/cryptic-stash/backend/ent"
	"github.com/jonboulle/clockwork"
)

func NewDummyUser(counter int, dbClient *ent.Client, ctx context.Context, clock clockwork.Clock) *ent.User {
	return dbClient.User.Create().
		SetUsername(fmt.Sprintf("user%v", counter)).
		SetSessionsValidFrom(clock.Now()).
		SetContent([]byte{1}).
		SetFileName("file.zip").
		SetMime("application/zip").
		SetNonce([]byte{1}).
		SetKeySalt([]byte{1}).
		SetHashTime(0).SetHashMemory(0).SetHashThreads(0).
		SaveX(ctx)
}
