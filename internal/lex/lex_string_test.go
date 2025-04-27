package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		input := "'hello world'"
		want := "hello world"
		cursor := Cursor{}

		got, newCursor, isValid := lexString(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("empty string", func(t *testing.T) {
		input := "''"
		cursor := Cursor{}

		got, newCursor, isValid := lexString(input, cursor)

		require.True(t, isValid)
		require.Equal(t, "", got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("unclosed string", func(t *testing.T) {
		input := "'hello world"
		cursor := Cursor{}

		_, _, isValid := lexString(input, cursor)

		require.False(t, isValid)
	})

	t.Run("invalid string", func(t *testing.T) {
		input := "hello world"
		cursor := Cursor{}

		_, _, isValid := lexString(input, cursor)

		require.False(t, isValid)
	})
}
