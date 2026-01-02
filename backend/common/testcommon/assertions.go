package testcommon

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NicoClack/cryptic-stash/backend/common"
	"github.com/stretchr/testify/require"
)

func AssertJSONEqual(t *testing.T, expected any, actual any, messagePrefix string) {
	messagePrefix += ": testcommon.AssertJSONEqual"
	expectedJSON, stdErr := json.Marshal(expected)
	require.NoError(t, stdErr, fmt.Sprintf("%v: marshalling expectedJSON shouldn't error", messagePrefix))
	actualJSON, stdErr := json.Marshal(actual)
	require.NoError(t, stdErr, fmt.Sprintf("%v: marshalling actualJSON shouldn't error", messagePrefix))

	require.Equal(t, string(expectedJSON), string(actualJSON), fmt.Sprintf("%v: JSON should match", messagePrefix))
}

func AssertJSONResponse(
	t *testing.T, respRecorder *httptest.ResponseRecorder,
	expectedStatus int,
	expectedPtr any,
) {
	if respRecorder.Code != expectedStatus {
		t.Fatalf(
			"expected HTTP status %v but got %v. response body:\n%v",
			expectedStatus, respRecorder.Code,
			respRecorder.Body.String(),
		)
	}
	expectedJSON, stdErr := json.Marshal(expectedPtr)
	require.NoError(t, stdErr)

	require.Equal(t, string(expectedJSON), respRecorder.Body.String())
}

func CallWithTimeout(t *testing.T, callback func(), timeout time.Duration) {
	select {
	case <-common.NewCallbackChannel(callback):
	case <-time.After(timeout):
		t.Fatalf("Function call timed out after %v", timeout)
	}
}
func AssertNoOp(t *testing.T, callback func()) {
	select {
	case <-common.NewCallbackChannel(callback):
	case <-time.After(5 * time.Millisecond):
		t.Fatalf("Expected no-op, but callback blocked")
	}
}
