package testcommon

import (
	"encoding/json"
	"fmt"
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
