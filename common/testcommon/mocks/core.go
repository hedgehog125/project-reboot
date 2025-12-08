package mocks

import (
	"context"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/ent"
)

type EmptyCoreService struct{}

func NewEmptyCoreService() *EmptyCoreService {
	return &EmptyCoreService{}
}

func (m *EmptyCoreService) RotateAdminCode() {
}
func (m *EmptyCoreService) CheckAdminCode(givenCode string) bool {
	return false
}
func (m *EmptyCoreService) RandomAuthCode() []byte {
	return []byte{}
}
func (m *EmptyCoreService) SendActiveSessionReminders(ctx context.Context) common.WrappedError {
	return nil
}
func (m *EmptyCoreService) DeleteExpiredSessions(ctx context.Context) common.WrappedError {
	return nil
}
func (m *EmptyCoreService) InvalidateUserSessions(userID int, ctx context.Context) common.WrappedError {
	return nil
}
func (m *EmptyCoreService) IsUserSufficientlyNotified(sessionOb *ent.Session) bool {
	return false
}
func (m *EmptyCoreService) IsUserLocked(userOb *ent.User) bool {
	return false
}
func (m *EmptyCoreService) Encrypt(data []byte, encryptionKey []byte) ([]byte, []byte, common.WrappedError) {
	return []byte{}, []byte{}, nil
}
func (m *EmptyCoreService) Decrypt(encrypted []byte, encryptionKey []byte, nonce []byte) ([]byte, common.WrappedError) {
	return []byte{}, nil
}
func (m *EmptyCoreService) GenerateSalt() []byte {
	return []byte{}
}
func (m *EmptyCoreService) HashPassword(password string, salt []byte, settings *common.PasswordHashSettings) []byte {
	return []byte{}
}
