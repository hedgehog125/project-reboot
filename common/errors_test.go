package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const errTypeTest = "test category [general]"

func TestGetSuccessfulActionIDs_returnsCorrectIDs(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		actionIDs []string
		errs      []*ErrWithStrId
		expected  []string
	}{
		{
			"empty",
			[]string{},
			[]*ErrWithStrId{},
			[]string{},
		},
		{
			"no errors",
			[]string{"action1", "action2", "action3"},
			[]*ErrWithStrId{},
			[]string{"action1", "action2", "action3"},
		},
		{
			"1/3 errors",
			[]string{"action1", "action2", "action3"},
			[]*ErrWithStrId{
				{
					Id:  "action1",
					Err: errors.New("error1"),
				},
			},
			[]string{"action2", "action3"},
		},
		{
			"2/3 errors",
			[]string{"action1", "action2", "action3"},
			[]*ErrWithStrId{
				{
					Id:  "action1",
					Err: errors.New("error1"),
				},
				{
					Id:  "action3",
					Err: errors.New("error3"),
				},
			},
			[]string{"action2"},
		},
		{
			"all failed",
			[]string{"action1", "action2", "action3"},
			[]*ErrWithStrId{
				{
					Id:  "action2",
					Err: errors.New("error2"),
				},
				{
					Id:  "action1",
					Err: errors.New("error1"),
				},
				{
					Id:  "action3",
					Err: errors.New("error3"),
				},
			},
			[]string{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ids := GetSuccessfulActionIDs(testCase.actionIDs, testCase.errs)
			require.Equal(t, testCase.expected, ids)
		})
	}
}

func TestError_Error_returnsCorrectMessage(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error",
		errTypeTest,
	)
	wrappedSentinelErr := sentinelErr.AddCategory("test function")
	databaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	)
	wrappedDatabaseErr := databaseErr.AddCategory("create user")

	require.Equal(t, "test error", sentinelErr.Error())
	require.Equal(t, "test function error: test error", wrappedSentinelErr.Error())
	require.Equal(t, "database error: database connection failed. details: ...", databaseErr.Error())
	require.Equal(t, "database error: create user error: database connection failed. details: ...", wrappedDatabaseErr.Error())
}

func TestError_worksWithIs(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error, no details",
		errTypeTest,
	)
	wrappedSentinelErr := sentinelErr.AddCategory("test function")
	databaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	)

	require.ErrorIs(t, sentinelErr, sentinelErr)
	require.NotErrorIs(t, sentinelErr, databaseErr)
	require.ErrorIs(t, sentinelErr, sentinelErr.Err)
	require.NotErrorIs(t, sentinelErr.Err, sentinelErr) // Target is more specific than err

	require.NotSame(t, sentinelErr, wrappedSentinelErr)
	require.ErrorIs(t, wrappedSentinelErr, sentinelErr)
	require.NotErrorIs(t, wrappedSentinelErr, databaseErr)
	require.ErrorIs(t, wrappedSentinelErr, wrappedSentinelErr.Err)
	require.NotErrorIs(t, wrappedSentinelErr.Err, wrappedSentinelErr) // Target is more specific than err
}

func TestError_HasCategories(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error, no details",
		errTypeTest,
	)
	flatDatabaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	)
	detailedDatabaseErr := WrapErrorWithCategories(
		errors.New("duplicate key error. details: ..."),
		ErrTypeDatabase,
		"create user",
	)

	require.True(t, sentinelErr.HasCategories(errTypeTest))
	require.True(t, sentinelErr.HasCategories("*"))
	require.False(t, sentinelErr.HasCategories(ErrTypeDatabase))
	require.True(t, sentinelErr.HasCategories(errTypeTest, "test error, no details"))
	require.True(t, sentinelErr.HasCategories(errTypeTest, "*"))
	require.True(t, sentinelErr.HasCategories("*", "test error, no details"))
	require.False(t, sentinelErr.HasCategories(ErrTypeDatabase, "test error, no details"))
	require.False(t, sentinelErr.HasCategories(ErrTypeDatabase, "*"))

	require.True(t, flatDatabaseErr.HasCategories(ErrTypeDatabase))
	require.True(t, flatDatabaseErr.HasCategories("*"))
	require.False(t, flatDatabaseErr.HasCategories(errTypeTest))
	require.False(t, flatDatabaseErr.HasCategories(ErrTypeDatabase, "some other category"))
	require.False(t, flatDatabaseErr.HasCategories(ErrTypeDatabase, "*"))
	require.False(t, flatDatabaseErr.HasCategories("*", "some other category"))
	require.False(t, flatDatabaseErr.HasCategories("*", "*"))

	require.False(t, detailedDatabaseErr.HasCategories(ErrTypeDatabase))
	require.False(t, detailedDatabaseErr.HasCategories(errTypeTest))
	require.True(t, detailedDatabaseErr.HasCategories("create user", ErrTypeDatabase))
	require.True(t, detailedDatabaseErr.HasCategories("create user", "*"))
	require.False(t, detailedDatabaseErr.HasCategories("create user", errTypeTest))
	require.True(t, detailedDatabaseErr.HasCategories("*", ErrTypeDatabase))
	require.True(t, detailedDatabaseErr.HasCategories("*", "*"))
	require.False(t, detailedDatabaseErr.HasCategories("*", errTypeTest))
}

func TestError_Clone(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error, no details",
		errTypeTest,
	)
	copiedErr := sentinelErr.Clone()

	require.Equal(t, sentinelErr, copiedErr)
	require.NotSame(t, sentinelErr, copiedErr)

	copiedErr.AddCategory("new category")
	require.NotEqual(t, sentinelErr, copiedErr)
}
