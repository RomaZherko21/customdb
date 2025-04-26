package lexer

import (
	"custom-database/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCreateTableCommand(t *testing.T) {
	t.Run("valid create table command", func(t *testing.T) {
		input := "CREATE TABLE films (id INT, title TEXT, rating INT, is_active BOOLEAN, release_date TEXT);"
		want := model.Table{
			TableName: "films",
			Columns: []model.Column{
				{Name: "id", Type: model.DataType("INT")},
				{Name: "title", Type: model.DataType("TEXT")},
				{Name: "rating", Type: model.DataType("INT")},
				{Name: "is_active", Type: model.DataType("BOOLEAN")},
				{Name: "release_date", Type: model.DataType("TEXT")},
			},
		}

		got, err := ParseCreateTableCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
		require.Equal(t, want.Columns, got.Columns)
	})

	t.Run("not enough arguments", func(t *testing.T) {
		input := "CREATE TABLE;"

		_, err := ParseCreateTableCommand(input)

		require.Error(t, err)
	})

	t.Run("no columns", func(t *testing.T) {
		input := "CREATE TABLE films;"

		_, err := ParseCreateTableCommand(input)

		require.Error(t, err)
	})

	t.Run("invalid column definition", func(t *testing.T) {
		input := "CREATE TABLE films (id);"

		_, err := ParseCreateTableCommand(input)

		require.Error(t, err)
	})
}

func TestExtractColumnsWithTypes(t *testing.T) {
	t.Run("valid columns with types", func(t *testing.T) {
		input := "CREATE TABLE films (id INT, title TEXT, rating INT, is_active BOOLEAN, release_date TEXT);"
		want := []model.Column{
			{Name: "id", Type: model.DataType("INT")},
			{Name: "title", Type: model.DataType("TEXT")},
			{Name: "rating", Type: model.DataType("INT")},
			{Name: "is_active", Type: model.DataType("BOOLEAN")},
			{Name: "release_date", Type: model.DataType("TEXT")},
		}

		got, err := extractColumnsWithTypes(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("no columns", func(t *testing.T) {
		input := "CREATE TABLE films;"

		_, err := extractColumnsWithTypes(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := extractColumnsWithTypes(input)

		require.Error(t, err)
	})

	t.Run("invalid format", func(t *testing.T) {
		input := "CREATE TABLE films id INT;"

		_, err := extractColumnsWithTypes(input)

		require.Error(t, err)
	})

	t.Run("invalid column definition", func(t *testing.T) {
		input := "CREATE TABLE films (id);"

		_, err := extractColumnsWithTypes(input)

		require.Error(t, err)
	})
}
