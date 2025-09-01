package testcommon

import (
	"context"
	"fmt"

	"github.com/hedgehog125/project-reboot/ent"
)

func NewDummyUser(counter int, dbClient *ent.Client, ctx context.Context) *ent.User {
	return dbClient.User.Create().SetUsername(fmt.Sprintf("user%v", counter)).
		SetContent([]byte{1}).SetFileName("file.zip").SetMime("application/zip").
		SetNonce([]byte{1}).SetKeySalt([]byte{1}).
		SetHashTime(0).SetHashMemory(0).SetHashThreads(0).
		SaveX(ctx)
}
