package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("echo test", func(t *testing.T) {
		code := RunCmd([]string{"echo", "1"}, Environment{})
		require.Equal(t, 0, code)
	})
}
