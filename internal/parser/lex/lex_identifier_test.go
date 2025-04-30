package lex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexIdentifier(t *testing.T) {
	t.Run("valid identifier", func(t *testing.T) {
		input := "table_name"
		want := "table_name"
		cursor := Cursor{}

		got, newCursor, isValid := lexIdentifier(input, cursor)

		require.True(t, isValid)
		require.Equal(t, want, got.Value)
		require.NotEqual(t, cursor, newCursor)
	})

	t.Run("invalid identifier", func(t *testing.T) {
		input := "123table"
		cursor := Cursor{}

		_, _, isValid := lexIdentifier(input, cursor)

		require.False(t, isValid)
	})

}
