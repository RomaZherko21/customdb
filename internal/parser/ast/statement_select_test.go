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
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "is_active"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "registered_at"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(10), cursor)
		require.Len(t, result.SelectedColumns, 4)
		require.Equal(t, "id", result.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.SelectedColumns[1].Literal.Value)
		require.Equal(t, "is_active", result.SelectedColumns[2].Literal.Value)
		require.Equal(t, "registered_at", result.SelectedColumns[3].Literal.Value)
		require.Equal(t, "users", result.From.Value)
	})

	t.Run("valid SELECT statement with FROM with limit and offset", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "is_active"},
			{Kind: lex.KeywordToken, Value: "from"},
			{Kind: lex.IdentifierToken, Value: "users"},
			{Kind: lex.KeywordToken, Value: "limit"},
			{Kind: lex.NumericToken, Value: "10"},
			{Kind: lex.KeywordToken, Value: "offset"},
			{Kind: lex.NumericToken, Value: "5"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, cursor, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, uint(12), cursor)
		require.Equal(t, "users", result.From.Value)
		require.Equal(t, 10, result.Limit)
		require.Equal(t, 5, result.Offset)
		require.Len(t, result.SelectedColumns, 3)
		require.Equal(t, "id", result.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.SelectedColumns[1].Literal.Value)
		require.Equal(t, "is_active", result.SelectedColumns[2].Literal.Value)
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

	t.Run("valid SELECT statement with complex WHERE", func(t *testing.T) {
		tokens := []*lex.Token{
			{Kind: lex.KeywordToken, Value: "select"},
			{Kind: lex.IdentifierToken, Value: "id"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "name"},
			{Kind: lex.SymbolToken, Value: ","},
			{Kind: lex.IdentifierToken, Value: "is_active"},
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
			{Kind: lex.LogicalOperatorToken, Value: "and"},
			{Kind: lex.IdentifierToken, Value: "is_active"},
			{Kind: lex.MathOperatorToken, Value: "="},
			{Kind: lex.BooleanToken, Value: "true"},
			{Kind: lex.SymbolToken, Value: ";"},
		}

		result, _, ok := parseSelectStatement(tokens, 0)

		require.True(t, ok)
		require.Equal(t, "users", result.From.Value)

		require.Len(t, result.SelectedColumns, 3)
		require.Equal(t, "id", result.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.SelectedColumns[1].Literal.Value)
		require.Equal(t, "is_active", result.SelectedColumns[2].Literal.Value)

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
		require.Equal(t, "and", result.Where.Right.Token.Value)

		require.Equal(t, "=", result.Where.Right.Left.Token.Value)
		require.Equal(t, "name", result.Where.Right.Left.Left.Token.Value)
		require.Equal(t, "'John'", result.Where.Right.Left.Right.Token.Value)

		// Проверяем правую часть (is_active = true)
		require.Equal(t, "=", result.Where.Right.Right.Token.Value)
		require.Equal(t, "is_active", result.Where.Right.Right.Left.Token.Value)
		require.Equal(t, "true", result.Where.Right.Right.Right.Token.Value)
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
