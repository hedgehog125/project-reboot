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
	).CommonError()
	wrappedSentinelErr := sentinelErr.AddCategory("test function")
	databaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	).CommonError()
	wrappedDatabaseErr := databaseErr.AddCategory("create user")
	packagedDatabaseErr := databaseErr.AddCategory("auth [package]")

	require.Equal(t, "test error", sentinelErr.Error())
	require.Equal(t, "test function error: test error", wrappedSentinelErr.Error())
	require.Equal(t, "database [general] error: database connection failed. details: ...", databaseErr.Error())
	require.Equal(t, "create user error: database [general] error: database connection failed. details: ...", wrappedDatabaseErr.Error())

	require.Equal(
		t,
		"auth [package] error: database [general] error: database connection failed. details: ...",
		packagedDatabaseErr.Error(),
	)
	packagedDatabaseErr = packagedDatabaseErr.AddCategory("create user")
	require.Equal(
		t,
		"auth [package] error: create user error: database [general] error: database connection failed. details: ...",
		packagedDatabaseErr.Error(),
	)

	repackagedDatabaseErr := packagedDatabaseErr.AddCategory("auth abstraction [package]")
	require.Equal(
		t,
		"auth abstraction [package] error: auth [package] error: create user error: database [general] error: database connection failed. details: ...",
		repackagedDatabaseErr.Error(),
	)
	repackagedDatabaseErr = repackagedDatabaseErr.AddCategory("abstraction function")
	require.Equal(
		t,
		"auth abstraction [package] error: abstraction function error: auth [package] error: create user error: database [general] error: database connection failed. details: ...",
		repackagedDatabaseErr.Error(),
	)
}

func TestError_worksWithIs(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error, no details",
		errTypeTest,
	).CommonError()
	wrappedSentinelErr := sentinelErr.AddCategory("test function")
	databaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	).CommonError()

	require.ErrorIs(t, sentinelErr, sentinelErr)
	require.NotErrorIs(t, sentinelErr, databaseErr)
	require.ErrorIs(t, sentinelErr, sentinelErr.err)
	require.NotErrorIs(t, sentinelErr.err, sentinelErr) // Target is more specific than err

	require.NotSame(t, sentinelErr, wrappedSentinelErr)
	require.ErrorIs(t, wrappedSentinelErr, sentinelErr)
	require.NotErrorIs(t, wrappedSentinelErr, databaseErr)
	require.ErrorIs(t, wrappedSentinelErr, wrappedSentinelErr.err)
	require.NotErrorIs(t, wrappedSentinelErr.err, wrappedSentinelErr) // Target is more specific than err
}

func TestError_HasCategories(t *testing.T) {
	t.Parallel()
	sentinelErr := NewErrorWithCategories(
		"test error, no details",
		errTypeTest,
	).CommonError()
	flatDatabaseErr := WrapErrorWithCategories(
		errors.New("database connection failed. details: ..."),
		ErrTypeDatabase,
	).CommonError()
	detailedDatabaseErr := WrapErrorWithCategories(
		errors.New("duplicate key error. details: ..."),
		ErrTypeDatabase,
		"create user",
	).CommonError()

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

	copiedErr = copiedErr.AddCategory("new category")
	require.NotEqual(t, sentinelErr, copiedErr)
}

func TestErrorWrapper(t *testing.T) {
	t.Parallel()
	databaseErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
	)
	createUserErrNoPackageWrapper := NewErrorWrapper(
		ErrTypeDatabase,
		"create user",
	)
	createUserErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
		"create user",
		"auth [package]",
	)

	rootError := errors.New("duplicate key error. details: ...")
	require.Equal(
		t,
		WrapErrorWithCategories(rootError, ErrTypeDatabase),
		databaseErrWrapper.Wrap(rootError),
	)
	require.Equal(
		t,
		WrapErrorWithCategories(rootError, ErrTypeDatabase, "create user"),
		createUserErrNoPackageWrapper.Wrap(rootError),
	)
	require.Equal(
		t,
		WrapErrorWithCategories(rootError, ErrTypeDatabase, "create user", "auth [package]"),
		createUserErrWrapper.Wrap(rootError),
	)
}

func TestErrorWrapper_removesDuplicatePackages(t *testing.T) {
	t.Parallel()
	errWrapperDbPackageA := NewErrorWrapper(
		"package A [package]", ErrTypeDatabase,
	)
	errWrapperCreateUserPackageA := NewErrorWrapper(
		"package A [package]", "create user",
	)
	errWrapperCreateUserPackageB := NewErrorWrapper(
		"package B [package]", "create team",
	)
	errWrapperSomethingPackageA := NewErrorWrapper(
		"package A [package]", "some category that is back in package A again",
	)

	rootError := errors.New("duplicate key error. details: ...")
	require.Equal(
		t,
		WrapErrorWithCategories(rootError, "package A [package]", ErrTypeDatabase),
		errWrapperDbPackageA.Wrap(rootError),
	)
	require.Equal(
		t,
		// The package should only appear once
		WrapErrorWithCategories(rootError, "package A [package]", ErrTypeDatabase, "create user"),
		errWrapperCreateUserPackageA.Wrap(errWrapperDbPackageA.Wrap(rootError)),
	)
	require.Equal(
		t,
		WrapErrorWithCategories(
			rootError, "package A [package]", ErrTypeDatabase, "create user",
			"package B [package]", "create team",
		),
		errWrapperCreateUserPackageB.Wrap(
			errWrapperCreateUserPackageA.Wrap(errWrapperDbPackageA.Wrap(rootError)),
		),
	)
	require.Equal(
		t,
		WrapErrorWithCategories(
			rootError, "package A [package]", ErrTypeDatabase, "create user",
			"package B [package]", "create team",
			"package A [package]", "some category that is back in package A again",
		),
		errWrapperSomethingPackageA.Wrap(
			errWrapperCreateUserPackageB.Wrap(
				errWrapperCreateUserPackageA.Wrap(errWrapperDbPackageA.Wrap(rootError)),
			),
		),
	)
}

func TestErrorWrapper_canAddPackageToPackagelessError(t *testing.T) {
	t.Parallel()
	commonErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
	)

	rootError := errors.New("duplicate key error. details: ...")
	wrappedError := commonErrWrapper.Wrap(rootError).CommonError().AddCategory("users [package]")
	require.Equal(
		t,
		[]string{ErrTypeDatabase, "users [package]"},
		wrappedError.Categories,
	)
	require.Equal(
		t,
		WrapErrorWithCategories(rootError, "users [package]", ErrTypeDatabase),
		wrappedError,
	)
}

func TestErrorWrapper_HasWrapped(t *testing.T) {
	t.Parallel()
	// TODO: wrap by passing the error through each wrapper
	// TODO: check this checks the packages properly
	commonDatabaseErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
	)
	authDatabaseErrWrapper := NewErrorWrapper(
		"auth [package]",
		ErrTypeDatabase,
	)
	createUserErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
		"auth [package]",
		"create user",
	)
	createUserAbstractionErrWrapper := NewErrorWrapper(
		ErrTypeDatabase,
		"auth [package]",
		"create user",
		"auth abstraction [package]",
		"abstraction function",
	)
	authPackageWrapper := NewErrorWrapper(
		"auth [package]",
	)

	rootError := errors.New("duplicate key error. details: ...")
	wrappedCommonDatabaseErr := commonDatabaseErrWrapper.Wrap(rootError).CommonError()
	wrappedAuthDatabaseErr := authDatabaseErrWrapper.Wrap(rootError).CommonError()
	wrappedCreateUserErr := createUserErrWrapper.Wrap(rootError).CommonError()
	wrappedCreateUserAbstractionErr := createUserAbstractionErrWrapper.Wrap(rootError).CommonError()

	require.False(t, createUserErrWrapper.HasWrapped(errors.New("generic error")))
	require.True(t, commonDatabaseErrWrapper.HasWrapped(wrappedCommonDatabaseErr))
	require.True(t, commonDatabaseErrWrapper.HasWrapped(wrappedAuthDatabaseErr)) // It compares the categories by value rather than tracking which wrappers were used
	require.True(t, authDatabaseErrWrapper.HasWrapped(wrappedAuthDatabaseErr))
	require.False(t, authDatabaseErrWrapper.HasWrapped(wrappedCommonDatabaseErr))
	require.False(t, createUserErrWrapper.HasWrapped(wrappedCommonDatabaseErr))
	require.False(t, createUserErrWrapper.HasWrapped(wrappedAuthDatabaseErr))
	require.True(t, commonDatabaseErrWrapper.HasWrapped(wrappedCreateUserErr))
	require.True(t, authDatabaseErrWrapper.HasWrapped(wrappedCreateUserErr))

	require.True(t, createUserErrWrapper.HasWrapped(wrappedCreateUserErr))
	require.False(t, createUserAbstractionErrWrapper.HasWrapped(wrappedCreateUserErr))
	require.True(t, createUserErrWrapper.HasWrapped(
		wrappedCreateUserErr.
			AddCategory("auth abstraction [package]").
			AddCategory("abstraction function"),
	))
	require.True(t, createUserErrWrapper.HasWrapped(
		wrappedCreateUserAbstractionErr,
	))

	require.True(t, authPackageWrapper.HasWrapped(
		wrappedCreateUserAbstractionErr,
	))
}
