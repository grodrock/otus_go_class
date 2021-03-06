package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("file not exist check", func(t *testing.T) {
		err := Copy("NotExistFile", "out", 0, 0)
		require.Errorf(t, os.ErrNotExist, err.Error())
	})

	t.Run("offset exceeds FileSize", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out", 1<<20, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("test copy big file with pb", func(t *testing.T) {
		t.Skip("File not present in repo")
		err := Copy("testdata/big_file.txt", "out", 0, 0)
		require.NoError(t, err)
	})
}
