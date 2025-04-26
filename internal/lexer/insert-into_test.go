package lexer

import (
	"custom-database/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseInsertIntoCommand(t *testing.T) {
	t.Run("valid insert command", func(t *testing.T) {
		input := "INSERT INTO films (code, title, did, date_prod, kind) VALUES ('T_601', 'Yojimbo', 106, '1961-06-16', 'Drama');"
		want := model.Table{
			TableName: "films",
			Columns: []model.Column{
				{Name: "code"},
				{Name: "title"},
				{Name: "did"},
				{Name: "date_prod"},
				{Name: "kind"},
			},
			Rows: [][]interface{}{
				{"'T_601'", "'Yojimbo'", "106", "'1961-06-16'", "'Drama'"},
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
