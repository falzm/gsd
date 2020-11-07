package gsd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_Get(t *testing.T) {
	state := State{}
	state.Store("test", "blah")

	require.Equal(t, "blah", state.Get("test"))
	require.Nil(t, state.Get("lolnope"))
}
