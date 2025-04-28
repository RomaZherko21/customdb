package ast

import (
	"custom-database/internal/lex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSelectStatement(t *testing.T) {
	t.Run("valid SELECT statement with FROM", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
		}
		delimiter := lex.Token{Kind: lex.SymbolToken, Value: ";"}

		result, cursor, ok := parseSelectStatement(tokens, 0, delimiter)

		require.True(t, ok)
		require.Equal(t, uint(6), cursor)
		require.Len(t, result.Item, 2)
		require.Equal(t, "id", result.Item[0].Literal.Value)
		require.Equal(t, "name", result.Item[1].Literal.Value)
		require.Equal(t, "users", result.From.Value)
	})

	t.Run("invalid SELECT statement - missing SELECT keyword", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
		}
		delimiter := lex.Token{Kind: lex.SymbolToken, Value: ";"}

		result, cursor, ok := parseSelectStatement(tokens, 0, delimiter)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid SELECT statement - missing expressions", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
		}
		delimiter := lex.Token{Kind: lex.SymbolToken, Value: ";"}

		result, cursor, ok := parseSelectStatement(tokens, 0, delimiter)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid SELECT statement - missing table name after FROM", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "from"},
		}
		delimiter := lex.Token{Kind: lex.SymbolToken, Value: ";"}

		result, cursor, ok := parseSelectStatement(tokens, 0, delimiter)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})
}
