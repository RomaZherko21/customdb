package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexBoolean(t *testing.T) {
	t.Run("valid true", func(t *testing.T) {
		input := "true"
		want := "true"
		cursor := Cursor{}

		got, newCursor, isValid := lexBoolean(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("valid false", func(t *testing.T) {
		input := "false"
		want := "false"
		cursor := Cursor{}

		got, newCursor, isValid := lexBoolean(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("invalid boolean", func(t *testing.T) {
		input := "not a boolean"
		cursor := Cursor{}

		_, _, isValid := lexBoolean(input, cursor)

		require.False(t, isValid)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""
		cursor := Cursor{}

		_, _, isValid := lexBoolean(input, cursor)

		require.False(t, isValid)
	})
}
