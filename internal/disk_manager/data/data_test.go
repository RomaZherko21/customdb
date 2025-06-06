package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertDataRow(t *testing.T) {
	// Создаем временный файл для тестов
	// tempDir := t.TempDir()
	tempDir := "./test"
	tempFile := filepath.Join(tempDir, "test.db")

	t.Run("успешная вставка строки после инициализации", func(t *testing.T) {
		// Создаем файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис
		ds := &dataService{
			file:          file,
			metaDataSpace: 0,
		}

		// Записываем метаданные
		columns := []Column{
			{
				Name:       "id",
				Type:       TypeInt32,
				IsNullable: false,
			},
			{
				Name:       "name",
				Type:       TypeText,
				IsNullable: true,
			},
		}

		metaData, metaSize, err := ds.writeMetaData("test_table", columns)
		assert.NoError(t, err)
		assert.NotNil(t, metaData)
		assert.Greater(t, metaSize, 0)
		ds.meta = metaData

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Вставляем тестовую строку
		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
			{Value: "test", Type: TypeText, IsNull: false},
		}

		result, err := ds.insertDataRow(row)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint32(INITIAL_PAGE_ID), result.PageID)
		assert.Equal(t, uint16(1), result.SlotID)

		// Проверяем, что данные действительно записались
		parsedRow, err := ds.ParseDataRow(result.PageID, result.SlotID)
		assert.NoError(t, err)
		assert.NotNil(t, parsedRow)
		assert.Len(t, parsedRow, len(row))
		assert.Equal(t, int32(1), parsedRow[0].Value)
		assert.Equal(t, "test", parsedRow[1].Value)
	})

	t.Run("успешная вставка строки с null значениями", func(t *testing.T) {
		// Создаем файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис
		ds := &dataService{
			file:          file,
			metaDataSpace: 0,
		}

		// Записываем метаданные
		columns := []Column{
			{
				Name:       "id",
				Type:       TypeInt32,
				IsNullable: false,
			},
			{
				Name:       "name",
				Type:       TypeText,
				IsNullable: true,
			},
		}

		metaData, metaSize, err := ds.writeMetaData("test_table", columns)
		assert.NoError(t, err)
		assert.NotNil(t, metaData)
		assert.Greater(t, metaSize, 0)
		ds.meta = metaData

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Вставляем тестовую строку с null значением
		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
			{Value: nil, Type: TypeText, IsNull: true},
		}

		result, err := ds.insertDataRow(row)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint32(INITIAL_PAGE_ID), result.PageID)
		assert.Equal(t, uint16(1), result.SlotID)

		// Проверяем, что данные действительно записались
		parsedRow, err := ds.ParseDataRow(result.PageID, result.SlotID)
		assert.NoError(t, err)
		assert.NotNil(t, parsedRow)
		assert.Len(t, parsedRow, len(row))
		assert.Equal(t, int32(1), parsedRow[0].Value)
		assert.Nil(t, parsedRow[1].Value)
		assert.True(t, parsedRow[1].IsNull)
	})

	t.Run("ошибка при вставке в несуществующий файл", func(t *testing.T) {
		// Создаем сервис с несуществующим файлом
		ds := &dataService{
			file:          nil,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
		}

		result, err := ds.insertDataRow(row)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("ошибка при вставке в неинициализированную таблицу", func(t *testing.T) {
		// Создаем пустой файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис без метаданных
		ds := &dataService{
			file:          file,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
		}

		result, err := ds.insertDataRow(row)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("успешная вставка нескольких строк", func(t *testing.T) {
		// Создаем файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис
		ds := &dataService{
			file:          file,
			metaDataSpace: 0,
		}

		// Записываем метаданные
		columns := []Column{
			{
				Name:       "id",
				Type:       TypeInt32,
				IsNullable: false,
			},
			{
				Name:       "name",
				Type:       TypeText,
				IsNullable: true,
			},
		}

		metaData, metaSize, err := ds.writeMetaData("test_table", columns)
		assert.NoError(t, err)
		assert.NotNil(t, metaData)
		assert.Greater(t, metaSize, 0)
		ds.meta = metaData

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Вставляем несколько строк
		rows := [][]DataCell{
			{
				{Value: int32(1), Type: TypeInt32, IsNull: false},
				{Value: "test1", Type: TypeText, IsNull: false},
			},
			{
				{Value: int32(2), Type: TypeInt32, IsNull: false},
				{Value: "test2", Type: TypeText, IsNull: false},
			},
			{
				{Value: int32(3), Type: TypeInt32, IsNull: false},
				{Value: "test3", Type: TypeText, IsNull: false},
			},
		}

		for i, row := range rows {
			result, err := ds.insertDataRow(row)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, uint32(INITIAL_PAGE_ID), result.PageID)
			assert.Equal(t, uint16(i+1), result.SlotID)

			// Проверяем каждую вставленную строку
			parsedRow, err := ds.ParseDataRow(result.PageID, result.SlotID)
			assert.NoError(t, err)
			assert.NotNil(t, parsedRow)
			assert.Len(t, parsedRow, len(row))
			assert.Equal(t, row[0].Value, parsedRow[0].Value)
			assert.Equal(t, row[1].Value, parsedRow[1].Value)
		}
	})
}
