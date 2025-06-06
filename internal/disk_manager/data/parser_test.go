package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePageHeader(t *testing.T) {
	// Создаем временный файл для тестов
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.db")

	t.Run("успешное чтение заголовка страницы после инициализации", func(t *testing.T) {
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

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Читаем заголовок страницы
		header, err := ds.ParsePageHeader(1)
		assert.NoError(t, err)
		assert.NotNil(t, header)

		// Проверяем поля заголовка
		assert.Equal(t, uint32(INITIAL_PAGE_ID), header.PageId)
		assert.Equal(t, uint16(DATA_SPACE), header.FreeSpace)
		assert.Equal(t, uint16(0), header.SlotsAmount)
	})

	t.Run("ошибка при чтении несуществующего файла", func(t *testing.T) {
		// Создаем сервис с несуществующим файлом
		ds := &dataService{
			file:          nil,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		header, err := ds.ParsePageHeader(1)
		assert.Error(t, err)
		assert.Nil(t, header)
	})

	t.Run("ошибка при чтении пустого файла", func(t *testing.T) {
		// Создаем пустой файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис
		ds := &dataService{
			file:          file,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		header, err := ds.ParsePageHeader(1)
		assert.Error(t, err)
		assert.Nil(t, header)
	})
}

func TestParsePageSlots(t *testing.T) {
	// Создаем временный файл для тестов
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.db")

	t.Run("успешное чтение слотов после инициализации", func(t *testing.T) {
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

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Читаем слоты страницы
		slots, err := ds.ParsePageSlots(1)
		assert.NoError(t, err)
		assert.NotNil(t, slots)
		assert.Len(t, slots, 0)
	})

	t.Run("ошибка при чтении несуществующего файла", func(t *testing.T) {
		// Создаем сервис с несуществующим файлом
		ds := &dataService{
			file:          nil,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		slots, err := ds.ParsePageSlots(1)
		assert.Error(t, err)
		assert.Nil(t, slots)
	})

	t.Run("ошибка при чтении пустого файла", func(t *testing.T) {
		// Создаем пустой файл
		file, err := os.Create(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		// Создаем сервис
		ds := &dataService{
			file:          file,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		slots, err := ds.ParsePageSlots(1)
		assert.Error(t, err)
		assert.Nil(t, slots)
	})
}

func TestParseDataRow(t *testing.T) {
	// Создаем временный файл для тестов
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.db")

	t.Run("ошибка при чтении несуществующего файла", func(t *testing.T) {
		// Создаем сервис с несуществующим файлом
		ds := &dataService{
			file:          nil,
			metaDataSpace: PAGE_HEADER_SIZE,
		}

		row, err := ds.ParseDataRow(1, 1)
		assert.Error(t, err)
		assert.Nil(t, row)
	})

	t.Run("ошибка при чтении несуществующего слота", func(t *testing.T) {
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
		}

		metaData, metaSize, err := ds.writeMetaData("test_table", columns)
		assert.NoError(t, err)
		assert.NotNil(t, metaData)
		assert.Greater(t, metaSize, 0)
		ds.meta = metaData

		// Записываем начальную страницу
		err = ds.writeInitialPageData()
		assert.NoError(t, err)

		// Пытаемся прочитать несуществующий слот
		row, err := ds.ParseDataRow(1, 999)
		assert.Error(t, err)
		assert.Nil(t, row)
	})
}
