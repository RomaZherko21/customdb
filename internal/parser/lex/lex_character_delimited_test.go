package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexCharacterDelimited(t *testing.T) {
	t.Run("valid delimited string", func(t *testing.T) {
		input := "'hello world'"
		want := "hello world"
		cursor := Cursor{}

		got, newCursor, isValid := lexCharacterDelimited(input, cursor, '\'')

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("empty string", func(t *testing.T) {
		input := "''"
		cursor := Cursor{}

		got, newCursor, isValid := lexCharacterDelimited(input, cursor, '\'')

		require.True(t, isValid)
		require.Equal(t, "", got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("unclosed string", func(t *testing.T) {
		input := "'hello world"
		cursor := Cursor{}

		_, _, isValid := lexCharacterDelimited(input, cursor, '\'')

		require.False(t, isValid)
	})

	t.Run("wrong delimiter", func(t *testing.T) {
		input := "'hello world'"
		cursor := Cursor{}

		_, _, isValid := lexCharacterDelimited(input, cursor, '?')

		require.False(t, isValid)
	})

	t.Run("double delimiter", func(t *testing.T) {
		input := "'he''llo'"
		cursor := Cursor{}

		got, _, isValid := lexCharacterDelimited(input, cursor, '\'')

		require.True(t, isValid)
		require.Equal(t, "he''llo", got.Value)
	})
}
