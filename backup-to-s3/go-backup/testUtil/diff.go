package testUtil

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"testing"
)

func RequireNoDiff(t *testing.T, expected, actual any) {
	diff := cmp.Diff(expected, actual)
	require.Empty(t, diff)
}
