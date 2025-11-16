package testcommon

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

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
