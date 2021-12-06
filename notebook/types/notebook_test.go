package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SupportedBookTypes(t *testing.T) {
	t.Run("all ok", func(t *testing.T) {
		bookTypes := SupportedBookTypes()
		require.NotEmpty(t, bookTypes)
	})
}
