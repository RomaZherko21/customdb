package ast

import (
	"custom-database/internal/parser/lex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSelectStatement(t *testing.T) {

	t.Run("valid SELECT statement with *", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.SymbolToken, Value: "*"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(4), cursor)
		require.Len(t, result.SelectedColumns, 0)
		require.Equal(t, "users", result.From.Value)
	})

	t.Run("valid SELECT statement with FROM", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(6), cursor)
		require.Len(t, result.SelectedColumns, 2)
		require.Equal(t, "id", result.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.SelectedColumns[1].Literal.Value)
		require.Equal(t, "users", result.From.Value)
	})

	t.Run("valid SELECT statement with simple WHERE", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "where"},
			{Kind: lex.IdentifierToken, Value: "age"},
			{Kind: lex.MathOperatorToken, Value: "="},
			{Kind: lex.NumericToken, Value: "18"},
			{Kind: lex.SymbolToken, Value: ";"},
		}
		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(8), cursor)
		require.Len(t, result.SelectedColumns, 1)
		require.Equal(t, "users", result.From.Value)

		// Проверяем условие WHERE (age = 18)
		require.NotNil(t, result.Where)
		require.Equal(t, "=", result.Where.Token.Value)
		require.Equal(t, "age", result.Where.Left.Token.Value)
		require.Equal(t, "18", result.Where.Right.Token.Value)
	})

	t.Run("valid SELECT statement with WHERE", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "where"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.MathOperatorToken, Value: ">"},
			{Kind: lex.NumericToken, Value: "1"},
			{Kind: lex.LogicalOperatorToken, Value: "and"},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.MathOperatorToken, Value: "="},
			{Kind: lex.StringToken, Value: "'John'"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, _, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, "users", result.From.Value)

		require.Len(t, result.SelectedColumns, 2)
		require.Equal(t, "id", result.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.SelectedColumns[1].Literal.Value)

		// Проверяем дерево WHERE условий
		require.NotNil(t, result.Where)
		require.Equal(t, "and", result.Where.Token.Value)

		// Проверяем левую часть (id > 1)
		require.NotNil(t, result.Where.Left)
		require.Equal(t, ">", result.Where.Left.Token.Value)
		require.Equal(t, "id", result.Where.Left.Left.Token.Value)
		require.Equal(t, "1", result.Where.Left.Right.Token.Value)

		// Проверяем правую часть (name = 'John')
		require.NotNil(t, result.Where.Right)
		require.Equal(t, "=", result.Where.Right.Token.Value)
		require.Equal(t, "name", result.Where.Right.Left.Token.Value)
		require.Equal(t, "'John'", result.Where.Right.Right.Token.Value)
	})

	t.Run("invalid SELECT statement - missing SELECT keyword", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid SELECT statement - missing expressions", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

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

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})

	t.Run("invalid SELECT statement - incomplete WHERE clause", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "where"},
			{Kind: lex.IdentifierToken, Value: "age"},
			{Kind: lex.MathOperatorToken, Value: "="},
			{Kind: lex.NumericToken, Value: "10"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.False(t, ok)
		require.Equal(t, uint(0), cursor)
		require.Nil(t, result)
	})
}
