package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStringBetween(t *testing.T) {
	require.Equal(t, "tag", GetStringBetween("test [tag]", "[", "]"))
	require.Equal(t, "", GetStringBetween("test", "[", "]"))
	require.Equal(t, "", GetStringBetween("test [", "[", "]"))
	require.Equal(t, "", GetStringBetween("test []", "[", "]"))

	require.Equal(t, "42", GetStringBetween("the meaning of life is 42.", "the meaning of life is ", "."))
}
