package parser

import (
	"custom-database/internal/parser/ast"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("valid CREATE TABLE statement", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT);"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, ast.CreateTableKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].CreateTableStatement.Name.Value)
		require.Len(t, *result.Statements[0].CreateTableStatement.Cols, 2)
		require.Equal(t, "id", (*result.Statements[0].CreateTableStatement.Cols)[0].Name.Value)
		require.Equal(t, "int", (*result.Statements[0].CreateTableStatement.Cols)[0].Datatype.Value)
		require.Equal(t, "name", (*result.Statements[0].CreateTableStatement.Cols)[1].Name.Value)
		require.Equal(t, "text", (*result.Statements[0].CreateTableStatement.Cols)[1].Datatype.Value)
	})

	t.Run("valid INSERT statement", func(t *testing.T) {
		source := "INSERT INTO users VALUES (1, 'Phil');"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, ast.InsertKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].InsertStatement.Table.Value)
		require.Len(t, *result.Statements[0].InsertStatement.Values, 2)
		require.Equal(t, "1", (*result.Statements[0].InsertStatement.Values)[0].Literal.Value)
		require.Equal(t, "Phil", (*result.Statements[0].InsertStatement.Values)[1].Literal.Value)
	})

	t.Run("valid SELECT statement", func(t *testing.T) {
		source := "SELECT id, name FROM users;"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, ast.SelectKind, result.Statements[0].Kind)
		require.Len(t, result.Statements[0].SelectStatement.SelectedColumns, 2)
		require.Equal(t, "id", result.Statements[0].SelectStatement.SelectedColumns[0].Literal.Value)
		require.Equal(t, "name", result.Statements[0].SelectStatement.SelectedColumns[1].Literal.Value)
		require.Equal(t, "users", result.Statements[0].SelectStatement.From.Value)
	})

	t.Run("valid DROP TABLE statement", func(t *testing.T) {
		source := "DROP TABLE users;"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, ast.DropTableKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].DropTableStatement.Table.Value)
	})

	t.Run("valid multiple statements", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT); INSERT INTO users VALUES (1, 'Phil');"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 2)

		// Проверяем CREATE TABLE
		require.Equal(t, ast.CreateTableKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].CreateTableStatement.Name.Value)
		require.Len(t, *result.Statements[0].CreateTableStatement.Cols, 2)

		// Проверяем INSERT
		require.Equal(t, ast.InsertKind, result.Statements[1].Kind)
		require.Equal(t, "users", result.Statements[1].InsertStatement.Table.Value)
		require.Len(t, *result.Statements[1].InsertStatement.Values, 2)
	})

	t.Run("invalid statement - missing semicolon", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT)"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - unknown statement type", func(t *testing.T) {
		source := "UNKNOWN users;"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed CREATE TABLE", func(t *testing.T) {
		source := "CREATE TABLE users (id INT name TEXT);"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed INSERT", func(t *testing.T) {
		source := "INSERT INTO users VALUES 1, 'Phil');"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed SELECT", func(t *testing.T) {
		source := "SELECT id name FROM users;"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed DROP TABLE", func(t *testing.T) {
		source := "DROP TABLE;"
		parser := NewParser()

		result, err := parser.Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})
}
