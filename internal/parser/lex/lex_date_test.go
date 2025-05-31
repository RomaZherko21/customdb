package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		input := "'2024-03-20 15:30:45'"
		want := "2024-03-20 15:30:45"
		cursor := Cursor{}

		got, newCursor, isValid := lexDate(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("empty string", func(t *testing.T) {
		input := "''"
		cursor := Cursor{}

		got, newCursor, isValid := lexDate(input, cursor)

		require.False(t, isValid)
		require.Nil(t, got)
		require.Equal(t, cursor, newCursor)
	})

	t.Run("unclosed string", func(t *testing.T) {
		input := "'hello world"
		cursor := Cursor{}

		_, _, isValid := lexDate(input, cursor)

		require.False(t, isValid)
	})

	t.Run("invalid date format", func(t *testing.T) {
		input := "'2024-58-59 15:30:45'"
		cursor := Cursor{}

		_, _, isValid := lexDate(input, cursor)

		require.False(t, isValid)
	})

	t.Run("invalid date format", func(t *testing.T) {
		input := "'2024-03-20'"
		cursor := Cursor{}

		_, _, isValid := lexDate(input, cursor)

		require.False(t, isValid)
	})
}
