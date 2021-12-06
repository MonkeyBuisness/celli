package utils

import (
	"errors"
	"testing"
)

type testCloser struct {
	err error
}

// Close closes the test closer.
func (e testCloser) Close() error {
	return e.err
}

func Test_Close(t *testing.T) {
	t.Run("close error", func(t *testing.T) {
		Close(testCloser{
			err: errors.New("error"),
		})
	})
	t.Run("all ok", func(t *testing.T) {
		Close(testCloser{})
	})
}
