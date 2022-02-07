package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read env dir", func(t *testing.T) {
		envExpected := Environment{
			"BAR":   EnvValue{"bar", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"EMPTY": EnvValue{"", false},
			"UNSET": EnvValue{"", true},
		}
		env, err := ReadDir("testdata/env")
		require.ErrorIs(t, nil, err)
		for k, v := range envExpected {
			require.Equal(t, v, env[k], "%v wrong value", k)
		}
	})
}

func TestReadEnvFile(t *testing.T) {
	t.Run("not exist file", func(t *testing.T) {
		_, err := ReadEnvFile("testdata/env/not_here")
		require.Error(t, err)
	})

	t.Run("file with second line", func(t *testing.T) {
		env, err := ReadEnvFile("testdata/env/BAR")
		require.Equal(t, "bar", env.Value)
		require.Equal(t, false, env.NeedRemove)
		require.Equal(t, nil, err)
	})

	t.Run("simple hello", func(t *testing.T) {
		env, err := ReadEnvFile("testdata/env/HELLO")
		require.Equal(t, "\"hello\"", env.Value)
		require.Equal(t, false, env.NeedRemove)
		require.Equal(t, nil, err)
	})

	t.Run("file with x00 separator", func(t *testing.T) {
		env, err := ReadEnvFile("testdata/env/FOO")
		require.Equal(t, "   foo\nwith new line", env.Value)
		require.Equal(t, false, env.NeedRemove)
		require.Equal(t, nil, err)
	})

	t.Run("empty file with symbols", func(t *testing.T) {
		env, err := ReadEnvFile("testdata/env/EMPTY")
		require.Equal(t, "", env.Value)
		require.Equal(t, false, env.NeedRemove)
		require.Equal(t, nil, err)
	})

	t.Run("empty file with zero length", func(t *testing.T) {
		env, err := ReadEnvFile("testdata/env/UNSET")
		require.Equal(t, "", env.Value)
		require.Equal(t, true, env.NeedRemove)
		require.Equal(t, nil, err)
	})
}

func TestUpdateOsEnv(t *testing.T) {
	t.Run("update env test", func(t *testing.T) {
		env := Environment{
			"CUSTOM1": EnvValue{"CUSTOM1", false},
			"CUSTOM2": EnvValue{"CUSTOM2", false},
		}
		err := env.UpdateOsEnv()
		require.ErrorIs(t, nil, err)
		require.Equal(t, "CUSTOM1", os.Getenv("CUSTOM1"))
		require.Equal(t, "CUSTOM2", os.Getenv("CUSTOM2"))
	})

	t.Run("unset env val", func(t *testing.T) {
		env := Environment{
			"CUSTOM1": EnvValue{"CUSTOM1", true},
		}
		env.UpdateOsEnv()
		require.Equal(t, "", os.Getenv("CUSTOM1"))
	})
}
