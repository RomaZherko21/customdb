package ast

import (
	"custom-database/internal/parser/lex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCreateTableStatement(t *testing.T) {
	t.Run("valid CREATE TABLE statement", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "create"},
			{Kind: lex.KeywordToken, Value: "table"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.SymbolToken, Value: "("},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "int"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "text"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "is_active"},
			{Kind: lex.KeywordToken, Value: "boolean"},
			{Kind: lex.SymbolToken, Value: ")"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseCreateTableStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(13), cursor)
		require.Equal(t, "users", result.Name.Value)
		require.Len(t, *result.Cols, 3)
		require.Equal(t, "id", (*result.Cols)[0].Name.Value)
		require.Equal(t, "int", (*result.Cols)[0].Datatype.Value)
		require.Equal(t, "name", (*result.Cols)[1].Name.Value)
		require.Equal(t, "text", (*result.Cols)[1].Datatype.Value)
		require.Equal(t, "is_active", (*result.Cols)[2].Name.Value)
		require.Equal(t, "boolean", (*result.Cols)[2].Datatype.Value)
	})

	t.Run("invalid CREATE statement - missing TABLE keyword", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "create"},
			{Kind: lex.IdentifierToken, Value: "users"},
		}

		result, cursor, ok := parseCreateTableStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid CREATE TABLE statement - missing table name", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "create"},
			{Kind: lex.KeywordToken, Value: "table"},
		}

		result, cursor, ok := parseCreateTableStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid CREATE TABLE statement - missing left parenthesis", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "create"},
			{Kind: lex.KeywordToken, Value: "table"},
			{Kind: lex.IdentifierToken, Value: "users"},
		}

		result, cursor, ok := parseCreateTableStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})
}

func TestParseColumnDefinitions(t *testing.T) {
	t.Run("valid column definitions", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "int"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "text"},
			{Kind: lex.SymbolToken, Value: ")"},
		}
		endDelimiter := lex.Token{Kind: lex.SymbolToken, Value: ")"}

		cols, cursor, ok := parseColumnDefinitions(tokens, 0, endDelimiter)

		require.True(t, ok)
		require.Equal(t, uint(5), cursor)
		require.Len(t, *cols, 2)
		require.Equal(t, "id", (*cols)[0].Name.Value)
		require.Equal(t, "int", (*cols)[0].Datatype.Value)
		require.Equal(t, "name", (*cols)[1].Name.Value)
		require.Equal(t, "text", (*cols)[1].Datatype.Value)
	})

	t.Run("invalid column definition - missing column type", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
		}
		endDelimiter := lex.Token{Kind: lex.SymbolToken, Value: ")"}

		cols, cursor, ok := parseColumnDefinitions(tokens, 0, endDelimiter)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, cols)
	})

	t.Run("invalid column definition - missing comma between columns", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "int"},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "text"},
		}
		endDelimiter := lex.Token{Kind: lex.SymbolToken, Value: ")"}

		cols, cursor, ok := parseColumnDefinitions(tokens, 0, endDelimiter)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, cols)
	})
}
