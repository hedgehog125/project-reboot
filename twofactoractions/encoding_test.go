package twofactoractions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode_DoesNotMutateTypeMap(t *testing.T) {
	registry := NewRegistry(nil)

	type NoAction1 struct {
		Foo string `json:"foo"`
	}
	registry.RegisterAction(ActionDefinition[any]{
		ID:       "NO_ACTION",
		Version:  1,
		BodyType: NoAction1{},
	})
	fullID := GetVersionedType("NO_ACTION", 1)

	actionDef, ok := registry.actions[fullID]
	require.True(t, ok)
	require.Equal(t, NoAction1{}, actionDef.BodyType)

	registry.Encode(fullID, NoAction1{
		Foo: "bar",
	})

	actionDef, ok = registry.actions[fullID]
	require.True(t, ok)
	require.Equal(t, NoAction1{}, actionDef.BodyType)
}
