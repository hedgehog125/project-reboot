package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPathPattern(t *testing.T) {
	t.Parallel()
	pattern1 := []string{"**", "documents", "projects", "**", "assets", "**"}
	pattern2 := []string{"***", "documents", "projects", "**", "assets", "**"}

	require.True(t, CheckPathPattern(
		[]string{"home", "nico", "documents", "projects", "unity", "experiments", "cool-game", "assets", "hats", "coolHat.png"},
		pattern1,
	))
	require.True(t, CheckPathPattern(
		[]string{"home", "jeff", "documents", "projects", "unity", "less-cool-game", "assets", "coats", "coolCoat.png"},
		pattern1,
	))
	require.False(t, CheckPathPattern(
		[]string{"home", "bob", "documents", "something-between", "projects", "unity", "experiments", "cool-game", "assets", "hats", "coolHat.png"},
		pattern1,
	))
	require.True(t, CheckPathPattern(
		[]string{"home", "alice", "documents", "something-between", "documents", "projects", "unity", "experiments", "cool-game", "assets", "hats", "coolHat.png"},
		pattern1,
	))
	require.False(t, CheckPathPattern(
		[]string{"random", "unrelated", "path"},
		pattern1,
	))

	require.False(t, CheckPathPattern( // ** should match 1 or more, not 0
		[]string{"documents", "projects", "unity", "experiments", "assets", "hats", "coolHat.png"},
		pattern1,
	))
	require.False(t, CheckPathPattern( // ** should match 1 or more, not 0
		[]string{"home", "anna", "documents", "projects", "unity", "less-cool-game", "assets"},
		pattern1,
	))
	require.True(t, CheckPathPattern( // *** matches 0 or more
		[]string{"documents", "projects", "unity", "experiments", "assets", "hats", "coolHat.png"},
		pattern2,
	))
	require.True(t, CheckPathPattern( // *** matches 0 or more
		[]string{"home", "anna", "documents", "projects", "unity", "less-cool-game", "assets"},
		[]string{"***", "documents", "projects", "**", "assets", "***"},
	))
	require.True(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"***", "apple", "banana"},
	))
	require.True(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"apple", "***", "banana"},
	))
	require.True(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"apple", "banana", "***"},
	))
	require.True(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"***", "apple", "*"},
	))
	require.False(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"***", "*", "apple", "*"}, // Should be treated as "**", "apple", "*"
	))
	require.True(t, CheckPathPattern(
		[]string{"something", "apple", "banana"},
		[]string{"***", "*", "apple", "*"},
	))
	require.False(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"**", "**", "apple", "*"}, // Should require at least 2 items before "apple"
	))
	require.True(t, CheckPathPattern(
		[]string{"item1", "item2", "apple", "banana"},
		[]string{"**", "**", "apple", "*"},
	))
	require.False(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"**", "***", "apple", "*"}, // Should be treated as "**", "apple", "*"
	))
	require.True(t, CheckPathPattern(
		[]string{"something", "apple", "banana"},
		[]string{"**", "***", "apple", "*"},
	))
	require.False(t, CheckPathPattern(
		[]string{"apple", "banana"},
		[]string{"*", "***", "apple", "*"},
	))
	require.True(t, CheckPathPattern(
		[]string{"something", "apple", "banana"},
		[]string{"*", "***", "apple", "*"},
	))
	require.False(t, CheckPathPattern(
		[]string{},
		[]string{"***", "*"},
	))
	require.True(t, CheckPathPattern(
		[]string{"something", "else"},
		[]string{"***", "*"},
	))
}
