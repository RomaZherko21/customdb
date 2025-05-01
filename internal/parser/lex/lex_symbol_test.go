package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexSymbol(t *testing.T) {
	t.Run("valid symbol", func(t *testing.T) {
		input := ";"
		want := ";"
		cursor := Cursor{}

		got, newCursor, isValid := lexSymbol(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("invalid symbol", func(t *testing.T) {
		input := "not a symbol"
		cursor := Cursor{}

		_, _, isValid := lexSymbol(input, cursor)

		require.False(t, isValid)
	})

}
