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
		fc, err := NewFileConnection(true, filename, tempDir, []meta.Column{})
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
		if page.Header.PageSize != PAGE_SIZE {
			t.Errorf("Неверный PageSize: получили %d, ожидали %d", page.Header.PageSize, PAGE_SIZE)
		}
	})

	// t.Run("Попытка создать файл с существующим именем", func(t *testing.T) {
	// 	// Подготовка
	// 	filename := "existing_table"

	// 	// Создаем файл перед тестом
	// 	existingFilePath := filepath.Join(tempDir, filename+".data")
	// 	file, _ := os.Create(existingFilePath)
	// 	file.Close()

	// 	// Действие
	// 	fc, err := NewFileConnection(true, filename, tempDir)
	// 	if err == nil {
	// 		t.Fatalf("Не удалось создать файл: %v", err)
	// 	}
	// 	fc.Close()
	// })

	t.Run("Создание файла с нулевым pageID", func(t *testing.T) {
		// Подготовка
		pageID := uint32(1)
		filename := "zero_page"

		// Действие
		fc, err := NewFileConnection(true, filename, tempDir, []meta.Column{})
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
		if page.Header.PageSize != PAGE_SIZE {
			t.Errorf("Неверный PageSize: получили %d, ожидали %d", page.Header.PageSize, PAGE_SIZE)
		}
	})
}
