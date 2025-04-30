package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexKeyword(t *testing.T) {
	t.Run("valid keyword", func(t *testing.T) {
		input := "SELECT"
		want := "select"
		cursor := Cursor{}

		got, newCursor, isValid := lexKeyword(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("invalid keyword", func(t *testing.T) {
		input := "not a keyword"
		cursor := Cursor{}

		_, _, isValid := lexKeyword(input, cursor)

		require.False(t, isValid)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""
		cursor := Cursor{}

		_, _, isValid := lexKeyword(input, cursor)

		require.False(t, isValid)
	})
}
