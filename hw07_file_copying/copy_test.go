package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	destinationFile = "/tmp/test-1.txt"
	sourceFile      = "./testdata/input.txt"
)

func TestCopy(t *testing.T) {
	t.Run("returns ErrOffsetExceedsFileSize error when offset is bigger than file size", func(t *testing.T) {
		err := Copy("./testdata/input.txt", destinationFile, 7000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("returns ErrUnsupportedFile error when source file is not regular file", func(t *testing.T) {
		err := Copy("/dev/random", destinationFile, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("returns ErrUnsupportedFile error when source file is directory", func(t *testing.T) {
		err := Copy("/tmp", destinationFile, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("test", func(t *testing.T) {
		err := Copy(sourceFile, destinationFile, 0, 10)
		require.NoError(t, err)
	})
}
