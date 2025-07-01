package servercommon

import (
	"errors"
	"testing"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/stretchr/testify/require"
)

func TestErrorUnwrap(t *testing.T) {
	stdErr := errors.New("test error")
	commErr := common.WrapErrorWithCategories(stdErr)
	serverErr := NewError(commErr)
	require.Equal(t, commErr, errors.Unwrap(serverErr))
	require.Equal(t, stdErr, errors.Unwrap(errors.Unwrap(serverErr)))
}
