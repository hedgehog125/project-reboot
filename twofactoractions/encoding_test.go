package twofactoractions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode_DoesNotMutateTypeMap(t *testing.T) {
	actionType, ok := actionTypeMap["No_ACTION_1"]
	require.True(t, ok)
	require.Equal(t, NoAction1{}, actionType)

	Encode("No_ACTION_1", NoAction1{
		Foo: "bar",
	})

	actionType, ok = actionTypeMap["No_ACTION_1"]
	require.True(t, ok)
	require.Equal(t, NoAction1{}, actionType)
}
