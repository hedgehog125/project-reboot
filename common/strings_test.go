package common_test

import (
	"testing"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/stretchr/testify/require"
)

func TestGetStringBetween(t *testing.T) {
	require.Equal(t, "tag", common.GetStringBetween("test [tag]", "[", "]"))
	require.Equal(t, "", common.GetStringBetween("test", "[", "]"))
	require.Equal(t, "", common.GetStringBetween("test [", "[", "]"))
	require.Equal(t, "", common.GetStringBetween("test []", "[", "]"))

	require.Equal(t, "42", common.GetStringBetween("the meaning of life is 42.", "the meaning of life is ", "."))
}

func TestParseVersionedType_GivenValidVersionedType_ReturnsParsed(t *testing.T) {
	jobType, version, err := common.ParseVersionedType("test_job_type_2")
	require.Nil(t, err)
	require.Equal(t, "test_job_type", jobType)
	require.Equal(t, 2, version)
}

func TestParseVersionedType_GivenUnderscoreSuffix_ReturnsError(t *testing.T) {
	_, _, err := common.ParseVersionedType("test_job_type_")
	require.ErrorIs(t, err, common.ErrMalformedVersionedType) // Shouldn't panic
}
