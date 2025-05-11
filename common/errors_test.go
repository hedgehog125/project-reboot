package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
