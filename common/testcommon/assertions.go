package testcommon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertJSONEqual(t *testing.T, expected any, actual any) {
	expectedJSON, stdErr := json.Marshal(expected)
	require.NoError(t, stdErr)
	actualJSON, stdErr := json.Marshal(actual)
	require.NoError(t, stdErr)

	require.Equal(t, string(expectedJSON), string(actualJSON))
}
