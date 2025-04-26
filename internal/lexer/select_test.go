package lexer

import (
	"custom-database/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSelectCommand(t *testing.T) {
	t.Run("valid select command", func(t *testing.T) {
		input := "SELECT name, id FROM users;"
		want := model.Table{
			TableName: "users",
			Columns: []model.Column{
				{Name: "name"},
				{Name: "id"},
			},
		}

		got, err := ParseSelectCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
		require.Equal(t, want.Columns, got.Columns)
	})

	t.Run("select all columns", func(t *testing.T) {
		input := "SELECT * FROM users;"
		want := model.Table{
			TableName: "users",
			Columns:   []model.Column{},
		}

		got, err := ParseSelectCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
		require.Equal(t, want.Columns, got.Columns)
	})

	t.Run("not enough arguments", func(t *testing.T) {
		input := "SELECT;"

		_, err := ParseSelectCommand(input)

		require.Error(t, err)
	})

	t.Run("no table name", func(t *testing.T) {
		input := "SELECT name, id FROM;"

		_, err := ParseSelectCommand(input)

		require.Error(t, err)
	})

	t.Run("invalid column name", func(t *testing.T) {
		input := "SELECT name id FROM users;"

		_, err := ParseSelectCommand(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := ParseSelectCommand(input)

		require.Error(t, err)
	})
}

func TestExtractSelectColumns(t *testing.T) {
	t.Run("valid select columns", func(t *testing.T) {
		input := "SELECT name, id FROM users;"
		want := []model.Column{
			{Name: "name"},
			{Name: "id"},
		}

		got, err := extractSelectColumns(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("single column", func(t *testing.T) {
		input := "SELECT name FROM users;"
		want := []model.Column{
			{Name: "name"},
		}

		got, err := extractSelectColumns(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("multiple columns with spaces", func(t *testing.T) {
		input := "SELECT name, id ,   age,  email FROM users;"
		want := []model.Column{
			{Name: "name"},
			{Name: "id"},
			{Name: "age"},
			{Name: "email"},
		}

		got, err := extractSelectColumns(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("select all columns", func(t *testing.T) {
		input := "SELECT * FROM users;"
		want := []model.Column{}

		got, err := extractSelectColumns(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("no columns", func(t *testing.T) {
		input := "SELECT FROM users;"

		_, err := extractSelectColumns(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := extractSelectColumns(input)

		require.Error(t, err)
	})

	t.Run("invalid format", func(t *testing.T) {
		input := "SELECT name id FROM users;"

		_, err := extractSelectColumns(input)

		require.Error(t, err)
	})
}
