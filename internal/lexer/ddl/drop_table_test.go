package ddl

import (
	"custom-database/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDropTableCommand(t *testing.T) {
	t.Run("valid drop table command", func(t *testing.T) {
		input := "DROP TABLE films;"
		want := model.Table{
			TableName: "films",
		}

		got, err := ParseDropTableCommand(input)

		require.NoError(t, err)
		require.Equal(t, want.TableName, got.TableName)
	})

	t.Run("not enough arguments", func(t *testing.T) {
		input := "DROP TABLE;"

		_, err := ParseDropTableCommand(input)

		require.Error(t, err)
	})

	t.Run("empty input", func(t *testing.T) {
		input := ""

		_, err := ParseDropTableCommand(input)

		require.Error(t, err)
	})

	t.Run("invalid command format", func(t *testing.T) {
		input := "DROP films;"

		_, err := ParseDropTableCommand(input)

		require.Error(t, err)
	})

	t.Run("table name with spaces", func(t *testing.T) {
		input := "DROP TABLE my films;"

		_, err := ParseDropTableCommand(input)

		require.Error(t, err)
	})
}
