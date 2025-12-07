package mocks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type ShutdownService struct {
	called bool
	reason string
}

func NewShutdownService() *ShutdownService {
	return &ShutdownService{}
}
func (service *ShutdownService) Listen() {
	panic("not implemented")
}
func (service *ShutdownService) Shutdown(reason string) {
	service.called = true
	service.reason = reason
}
func (service *ShutdownService) AssertCalled(t *testing.T, expectedReason string) {
	require.True(t, service.called)
	require.Equal(t, expectedReason, service.reason)
}
func (service *ShutdownService) AssertNotCalled(t *testing.T) {
	require.False(t, service.called)
}
func (service *ShutdownService) Reset() {
	service.called = false
	service.reason = ""
}
