package testcommon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertJSONEqual(t *testing.T, expected any, actual any) {
	expectedJSON, stdErr := json.MarshalIndent(expected, "", "\t")
	require.NoError(t, stdErr)
	actualJSON, stdErr := json.MarshalIndent(actual, "", "\t")
	require.NoError(t, stdErr)

	require.JSONEq(t, string(expectedJSON), string(actualJSON))
}
