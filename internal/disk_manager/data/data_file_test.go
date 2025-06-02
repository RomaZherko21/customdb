package data

import (
	"custom-database/internal/disk_manager/meta"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDataFile(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	t.Run("Успешное создание файла", func(t *testing.T) {
		// Подготовка
		pageID := uint32(1)
		filename := "test_table"

		// Действие
		fc, err := NewFileConnection(&meta.MetaFile{
			Name:      "test_table",
			PageCount: 1,
			Columns:   []meta.Column{},
		}, tempDir, true)
		if err != nil {
			t.Fatalf("Не удалось создать файл: %v", err)
		}
		defer fc.Close()

		// Проверка
		if err != nil {
			t.Errorf("Неожиданная ошибка: %v", err)
		}

		// Проверяем файл
		filePath := filepath.Join(tempDir, filename+".data")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("Файл не был создан")
		}

		// Проверяем размер и содержимое
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("Не удалось открыть созданный файл: %v", err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			t.Fatalf("Не удалось получить информацию о файле: %v", err)
		}

		if fileInfo.Size() != int64(PAGE_SIZE) {
			t.Errorf("Неверный размер файла: получили %d, ожидали %d", fileInfo.Size(), PAGE_SIZE)
		}

		page := fc.deserializePage(pageID)
		if page.Header.PageId != pageID {
			t.Errorf("Неверный PageID: получили %d, ожидали %d", page.Header.PageId, pageID)
		}
		if page.Header.FreeSpace != DATA_SPACE {
			t.Errorf("Неверный PageSize: получили %d, ожидали %d", page.Header.FreeSpace, DATA_SPACE)
		}
	})

	t.Run("Создание файла с нулевым pageID", func(t *testing.T) {
		// Подготовка
		pageID := uint32(1)
		filename := "zero_page"

		// Действие
		fc, err := NewFileConnection(&meta.MetaFile{
			Name:      "zero_page",
			PageCount: 1,
			Columns:   []meta.Column{},
		}, tempDir, true)
		if err != nil {
			t.Fatalf("Не удалось создать файл: %v", err)
		}
		defer fc.Close()

		// Проверка
		if err != nil {
			t.Errorf("Неожиданная ошибка: %v", err)
		}

		// Проверяем файл
		filePath := filepath.Join(tempDir, filename+".data")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("Файл не был создан")
		}

		// Проверяем размер и содержимое
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatalf("Не удалось открыть созданный файл: %v", err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			t.Fatalf("Не удалось получить информацию о файле: %v", err)
		}

		if fileInfo.Size() != int64(PAGE_SIZE) {
			t.Errorf("Неверный размер файла: получили %d, ожидали %d", fileInfo.Size(), PAGE_SIZE)
		}

		page := fc.deserializePage(pageID)
		if page.Header.PageId != pageID {
			t.Errorf("Неверный PageID: получили %d, ожидали %d", page.Header.PageId, pageID)
		}
		if page.Header.FreeSpace != DATA_SPACE {
			t.Errorf("Неверный PageSize: получили %d, ожидали %d", page.Header.FreeSpace, DATA_SPACE)
		}
	})
}
