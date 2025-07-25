package jobscommon

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersionedType_givenValidVersionedType_returnsParsed(t *testing.T) {
	jobType, version, err := ParseVersionedType("test_job_type_2")
	require.Nil(t, err)
	require.Equal(t, "test_job_type", jobType)
	require.Equal(t, 2, version)
}

func TestParseVersionedType_givenUnderscoreSuffix_returnsError(t *testing.T) {
	_, _, err := ParseVersionedType("test_job_type_")
	require.ErrorIs(t, err, ErrMalformedVersionedType) // Shouldn't panic
}
