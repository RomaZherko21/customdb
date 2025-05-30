package ast

import (
	"custom-database/internal/parser/lex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInsertStatement(t *testing.T) {
	t.Run("valid INSERT statement", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.KeywordToken, Value: "into"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "values"},
			{Kind: lex.SymbolToken, Value: "("},
			{Kind: lex.NumericToken, Value: "1"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.StringToken, Value: "John"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.BooleanToken, Value: "true"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.BooleanToken, Value: "false"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.NullToken, Value: "null"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.DateToken, Value: "2024-03-20 15:30:45"},
			{Kind: lex.SymbolToken, Value: ")"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(17), cursor)
		require.Equal(t, "users", result.Table.Value)
		require.Len(t, *result.Values, 6)
		require.Equal(t, "1", (*result.Values)[0].Literal.Value)
		require.Equal(t, "John", (*result.Values)[1].Literal.Value)
		require.Equal(t, "true", (*result.Values)[2].Literal.Value)
		require.Equal(t, "false", (*result.Values)[3].Literal.Value)
		require.Equal(t, "null", (*result.Values)[4].Literal.Value)
		require.Equal(t, "2024-03-20 15:30:45", (*result.Values)[5].Literal.Value)
	})

	t.Run("invalid INSERT statement - missing INTO keyword", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.IdentifierToken, Value: "users"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid INSERT statement - missing table name", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.KeywordToken, Value: "into"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid INSERT statement - missing VALUES keyword", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.KeywordToken, Value: "into"},
			{Kind: lex.IdentifierToken, Value: "users"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid INSERT statement - missing left parenthesis", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.KeywordToken, Value: "into"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "values"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid INSERT statement - missing right parenthesis", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "insert"},
			{Kind: lex.KeywordToken, Value: "into"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "values"},
			{Kind: lex.SymbolToken, Value: "("},
			{Kind: lex.NumericToken, Value: "1"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.StringToken, Value: "John"},
		}

		result, cursor, ok := parseInsertStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})
}
