package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ErrorsIs(t *testing.T) {
	err := ErrReadMarkdownSource.New("test")
	require.True(t, errors.Is(err, ErrReadMarkdownSource))
	require.False(t, errors.Is(err, ErrMarshalCommentPayload))
}

func TestError_New(t *testing.T) {
	details := []string{"test", "test2", "test3"}
	err := ErrReadMarkdownSource.New(details...)
	require.Error(t, err)
	require.EqualError(t, err.base, ErrReadMarkdownSource.base.Error())
	require.ElementsMatch(t, err.details, details)
}

func TestError_Error(t *testing.T) {
	details := []string{"test", "test2", "test3"}
	err := ErrReadMarkdownSource.New(details...)
	require.Equal(t, fmt.Sprintf("%v: %s", err.base, strings.Join(details, "; ")), err.Error())
}

func TestError_Is(t *testing.T) {
	err := ErrReadMarkdownSource.New("test")
	require.True(t, err.Is(ErrReadMarkdownSource))
	require.False(t, err.Is(ErrMarshalCommentPayload))
}
