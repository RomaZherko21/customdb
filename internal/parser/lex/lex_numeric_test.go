package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexNumeric(t *testing.T) {
	t.Run("valid integer", func(t *testing.T) {
		input := "42"
		want := "42"
		cursor := Cursor{}

		got, newCursor, isValid := lexNumeric(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("valid float", func(t *testing.T) {
		input := "3.14"
		want := "3.14"
		cursor := Cursor{}

		got, newCursor, isValid := lexNumeric(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("invalid number", func(t *testing.T) {
		input := "not a number"
		cursor := Cursor{}

		_, _, isValid := lexNumeric(input, cursor)

		require.False(t, isValid)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""
		cursor := Cursor{}

		_, _, isValid := lexNumeric(input, cursor)

		require.False(t, isValid)
	})
}
