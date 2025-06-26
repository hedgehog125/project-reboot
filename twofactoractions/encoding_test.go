package twofactoractions

import (
	"errors"
	"testing"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/stretchr/testify/require"
)

func TestEncode_givenWrongType_returnsErr(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(nil)

	type NoAction1 struct {
		Foo string `binding:"required" json:"foo"`
	}
	registry.RegisterAction(&ActionDefinition{
		ID:       "NO_ACTION",
		Version:  1,
		BodyType: &NoAction1{},
	})

	newInvalidDataError := func(message string) *common.Error {
		return ErrWrapperInvalidData.Wrap(
			errors.New(message),
		).AddCategory(ErrTypeEncode)
	}
	_, commErr := registry.Encode("NO_ACTION_1", &struct{}{})
	require.Equal(t, commErr, newInvalidDataError("data type *struct {} isn't the expected type *twofactoractions.NoAction1"))
	_, commErr = registry.Encode("NO_ACTION_1", 42)
	require.Equal(t, commErr, newInvalidDataError("data type int isn't the expected type *twofactoractions.NoAction1"))

	_, commErr = registry.Encode("NO_ACTION_1", nil)
	require.Equal(t, commErr, newInvalidDataError("data type %!s(<nil>) isn't the expected type *twofactoractions.NoAction1")) // I guess this is better than panicking?
	var nilInterface any = nil
	_, commErr = registry.Encode("NO_ACTION_1", nilInterface)
	require.Equal(t, commErr, newInvalidDataError("data type %!s(<nil>) isn't the expected type *twofactoractions.NoAction1")) // I guess this is better than panicking?

	type SimilarAction struct {
		Foo string `binding:"required" json:"foo"`
	}
	type OtherSimilarAction struct {
		Foo string `binding:"required" json:"foo"`
		Bar int    `json:"bar"`
	}
	_, commErr = registry.Encode("NO_ACTION_1", &SimilarAction{
		Foo: "bar",
	})
	require.Equal(t, commErr, newInvalidDataError("data type *twofactoractions.SimilarAction isn't the expected type *twofactoractions.NoAction1"))
	_, commErr = registry.Encode("NO_ACTION_1", &OtherSimilarAction{
		Foo: "bar",
		Bar: 42,
	})
	require.Equal(t, commErr, newInvalidDataError("data type *twofactoractions.OtherSimilarAction isn't the expected type *twofactoractions.NoAction1"))

	_, commErr = registry.Encode("NO_ACTION_1", NoAction1{ // Note: this is a value, not a pointer
		Foo: "bar",
	})
	require.Equal(t, commErr, newInvalidDataError("data type twofactoractions.NoAction1 isn't the expected type *twofactoractions.NoAction1"))
}

func TestEncode_givenCorrectType_returnsJsonString(t *testing.T) {
	t.Parallel()
	registry := NewRegistry(nil)

	type NoAction1 struct {
		Foo string `binding:"required" json:"foo"`
	}
	registry.RegisterAction(&ActionDefinition{
		ID:       "NO_ACTION",
		Version:  1,
		BodyType: &NoAction1{},
	})

	encoded, commErr := registry.Encode("NO_ACTION_1", &NoAction1{
		Foo: "bar",
	})
	require.Nil(t, commErr)
	require.Equal(t, "{\"foo\":\"bar\"}", encoded)
}
