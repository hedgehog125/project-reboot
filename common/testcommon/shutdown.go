package testcommon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MockShutdownService struct {
	called bool
	reason string
}

func NewMockShutdownService() *MockShutdownService {
	return &MockShutdownService{}
}
func (service *MockShutdownService) Listen() {
	panic("not implemented")
}
func (service *MockShutdownService) Shutdown(reason string) {
	service.called = true
	service.reason = reason
}
func (service *MockShutdownService) AssertCalled(t *testing.T, expectedReason string) {
	require.True(t, service.called)
	require.Equal(t, expectedReason, service.reason)
}
func (service *MockShutdownService) AssertNotCalled(t *testing.T) {
	require.False(t, service.called)
}
func (service *MockShutdownService) Reset() {
	service.called = false
	service.reason = ""
}
