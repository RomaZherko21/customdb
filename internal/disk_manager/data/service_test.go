package data

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTableData(t *testing.T) {
	t.Run("Создание таблицы", func(t *testing.T) {
		// Создаем тестовую директорию
		err := os.MkdirAll("./test_data", 0755)
		assert.NoError(t, err)

		file, err := os.Create("./test_data/test_table.data")
		assert.NoError(t, err)

		// Создаем тестовые колонки
		columns := []Column{
			{Name: "id", Type: TypeInt32, IsNullable: false},
			{Name: "name", Type: TypeText, IsNullable: true},
		}

		_, err = InitTableData(file, "test_table", columns)
		assert.NoError(t, err)
	})

	t.Run("Вставка данных в таблицу", func(t *testing.T) {
		file, err := os.Open("./test_data/test_table.data")
		assert.NoError(t, err)

		ds, err := NewDataService(file, "test_table")
		assert.NoError(t, err)

		ds.InsertDataRow([]DataCell{
			{Value: int32(1), Type: TypeInt32},
			{Value: "test", Type: TypeText},
		})
		assert.NoError(t, err)

		ds.InsertDataRow([]DataCell{
			{Value: int32(2), Type: TypeInt32},
			{Value: "test2", Type: TypeText},
		})
		assert.NoError(t, err)
	})
}
