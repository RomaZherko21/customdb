package lexer

import (
	"custom-database/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractValue(t *testing.T) {
	t.Run("string value", func(t *testing.T) {
		input := "'hello'"
		want := "hello"

		got, err := extractValue(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("integer value", func(t *testing.T) {
		input := "42"
		want := 42

		got, err := extractValue(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("boolean value", func(t *testing.T) {
		input := "true"
		want := true

		got, err := extractValue(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("null value", func(t *testing.T) {
		input := "NULL"

		got, err := extractValue(input)

		require.NoError(t, err)
		require.Equal(t, nil, got)
	})

	t.Run("invalid value", func(t *testing.T) {
		input := "invalid"

		_, err := extractValue(input)

		require.Error(t, err)
	})
}

func TestParseInsertIntoCommand(t *testing.T) {
	t.Run("valid insert command with all types", func(t *testing.T) {
		input := "INSERT INTO films (id, title, rating, is_active, release_date) VALUES (1, 'The Matrix', 8, true, '1999-03-31');"
		want := model.Table{
			TableName: "films",
			Columns: []model.Column{
				{Name: "id"},
				{Name: "title"},
				{Name: "rating"},
				{Name: "is_active"},
				{Name: "release_date"},
			},
			Rows: [][]interface{}{
				{1, "The Matrix", 8, true, "1999-03-31"},
			},
		}

		got, err := ParseInsertIntoCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
		require.Equal(t, want.Columns, got.Columns)
		require.Equal(t, want.Rows, got.Rows)
	})

	t.Run("valid insert command with null values", func(t *testing.T) {
		input := "INSERT INTO films (id, title, rating) VALUES (1, 'The Matrix', NULL);"
		want := model.Table{
			TableName: "films",
			Columns: []model.Column{
				{Name: "id"},
				{Name: "title"},
				{Name: "rating"},
			},
			Rows: [][]interface{}{
				{1, "The Matrix", nil},
			},
		}

		got, err := ParseInsertIntoCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
		require.Equal(t, want.Columns, got.Columns)
		require.Equal(t, want.Rows, got.Rows)
	})

	t.Run("not enough arguments", func(t *testing.T) {
		input := "INSERT INTO;"

		_, err := ParseInsertIntoCommand(input)

		require.Error(t, err)
	})

	t.Run("no columns", func(t *testing.T) {
		input := "INSERT INTO films VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');"

		_, err := ParseInsertIntoCommand(input)

		require.Error(t, err)
	})

	t.Run("no values", func(t *testing.T) {
		input := "INSERT INTO films (code, title, did, date_prod, kind);"

		_, err := ParseInsertIntoCommand(input)

		require.Error(t, err)
	})

	t.Run("mismatched columns and values", func(t *testing.T) {
		input := "INSERT INTO films (id, title) VALUES (1);"

		_, err := ParseInsertIntoCommand(input)

		require.Error(t, err)
	})
}

func TestExtractColumns(t *testing.T) {
	t.Run("valid columns", func(t *testing.T) {
		input := "INSERT INTO films (code, title, did, date_prod, kind) VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');"
		want := []string{"code", "title", "did", "date_prod", "kind"}

		got, err := extractColumns(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("no columns", func(t *testing.T) {
		input := "INSERT INTO films VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');"

		_, err := extractColumns(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := extractColumns(input)

		require.Error(t, err)
	})
}

func TestExtractValues(t *testing.T) {
	t.Run("valid values", func(t *testing.T) {
		input := "INSERT INTO films (code, title, did, date_prod, kind) VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');"
		want := []string{"'T_601'", "'Yojimbo'", "106", "'1961-06-16'", "'Drama'"}

		got, err := extractValues(input)

		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("no values", func(t *testing.T) {
		input := "INSERT INTO films (code, title, did, date_prod, kind);"

		_, err := extractValues(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := extractValues(input)

		require.Error(t, err)
	})
}
