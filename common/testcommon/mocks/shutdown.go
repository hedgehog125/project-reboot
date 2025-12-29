package mocks

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type ShutdownService struct {
	called bool
	reason string
	mu     sync.Mutex
}

func NewShutdownService() *ShutdownService {
	return &ShutdownService{}
}
func (service *ShutdownService) Listen() {
	panic("not implemented")
}
func (service *ShutdownService) Shutdown(reason string) {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.called = true
	service.reason = reason
}
func (service *ShutdownService) AssertCalled(t *testing.T, expectedReason string) {
	service.mu.Lock()
	defer service.mu.Unlock()
	require.True(t, service.called)
	require.Equal(t, expectedReason, service.reason)
}
func (service *ShutdownService) AssertNotCalled(t *testing.T) {
	service.mu.Lock()
	defer service.mu.Unlock()
	require.False(t, service.called)
}
func (service *ShutdownService) Reset() {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.called = false
	service.reason = ""
}
