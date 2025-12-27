package servercommon_test

import (
	"errors"
	"testing"

	"github.com/NicoClack/cryptic-stash/common"
	"github.com/NicoClack/cryptic-stash/server/servercommon"
	"github.com/stretchr/testify/require"
)

func TestErrorUnwrap(t *testing.T) {
	t.Parallel()
	stdErr := errors.New("test error")
	wrappedErr := common.WrapErrorWithCategories(stdErr)
	serverErr := servercommon.NewError(wrappedErr)
	require.Equal(t, wrappedErr, errors.Unwrap(serverErr))
	require.Equal(t, stdErr, errors.Unwrap(errors.Unwrap(serverErr)))
}

func TestErrorImplementsWrappedError(t *testing.T) {
	t.Parallel()
	serverErr := servercommon.NewError(common.NewErrorWithCategories("test error", "testing [package]"))
	var _ common.WrappedError = serverErr
}

func TestErrorIsNotUnwrappedByWithRetries(t *testing.T) {
	t.Parallel()
	serverErr := servercommon.NewError(
		common.NewErrorWithCategories("not found", "testing [package]"),
	).SetStatus(404)

	wrappedErr := common.WithRetries(t.Context(), nil, func() error {
		return serverErr
	})
	serverErr = servercommon.NewError(wrappedErr)
	require.Equal(t, 404, serverErr.Status())
}
