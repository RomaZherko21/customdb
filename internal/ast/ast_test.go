package ast

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("valid CREATE TABLE statement", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT);"

		result, err := Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, CreateTableKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].CreateTableStatement.name.Value)
		require.Len(t, *result.Statements[0].CreateTableStatement.cols, 2)
		require.Equal(t, "id", (*result.Statements[0].CreateTableStatement.cols)[0].name.Value)
		require.Equal(t, "int", (*result.Statements[0].CreateTableStatement.cols)[0].datatype.Value)
		require.Equal(t, "name", (*result.Statements[0].CreateTableStatement.cols)[1].name.Value)
		require.Equal(t, "text", (*result.Statements[0].CreateTableStatement.cols)[1].datatype.Value)
	})

	t.Run("valid INSERT statement", func(t *testing.T) {
		source := "INSERT INTO users VALUES (1, 'Phil');"

		result, err := Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, InsertKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].InsertStatement.table.Value)
		require.Len(t, *result.Statements[0].InsertStatement.values, 2)
		require.Equal(t, "1", (*result.Statements[0].InsertStatement.values)[0].literal.Value)
		require.Equal(t, "Phil", (*result.Statements[0].InsertStatement.values)[1].literal.Value)
	})

	t.Run("valid SELECT statement", func(t *testing.T) {
		source := "SELECT id, name FROM users;"

		result, err := Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 1)
		require.Equal(t, SelectKind, result.Statements[0].Kind)
		require.Len(t, result.Statements[0].SelectStatement.item, 2)
		require.Equal(t, "id", result.Statements[0].SelectStatement.item[0].literal.Value)
		require.Equal(t, "name", result.Statements[0].SelectStatement.item[1].literal.Value)
		require.Equal(t, "users", result.Statements[0].SelectStatement.from.Value)
	})

	t.Run("valid multiple statements", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT); INSERT INTO users VALUES (1, 'Phil');"

		result, err := Parse(source)

		require.NoError(t, err)
		require.Len(t, result.Statements, 2)

		// Проверяем CREATE TABLE
		require.Equal(t, CreateTableKind, result.Statements[0].Kind)
		require.Equal(t, "users", result.Statements[0].CreateTableStatement.name.Value)
		require.Len(t, *result.Statements[0].CreateTableStatement.cols, 2)

		// Проверяем INSERT
		require.Equal(t, InsertKind, result.Statements[1].Kind)
		require.Equal(t, "users", result.Statements[1].InsertStatement.table.Value)
		require.Len(t, *result.Statements[1].InsertStatement.values, 2)
	})

	t.Run("invalid statement - missing semicolon", func(t *testing.T) {
		source := "CREATE TABLE users (id INT, name TEXT)"

		result, err := Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - unknown statement type", func(t *testing.T) {
		source := "UNKNOWN users;"

		result, err := Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed CREATE TABLE", func(t *testing.T) {
		source := "CREATE TABLE users (id INT name TEXT);"

		result, err := Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed INSERT", func(t *testing.T) {
		source := "INSERT INTO users VALUES 1, 'Phil');"

		result, err := Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("invalid statement - malformed SELECT", func(t *testing.T) {
		source := "SELECT id name FROM users;"

		result, err := Parse(source)

		require.Error(t, err)
		require.Nil(t, result)
	})
}
