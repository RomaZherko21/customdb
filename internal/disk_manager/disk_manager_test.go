package disk_manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "Create file successfully",
			filename: "test.db",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем после теста
			defer func() {
				if !tt.wantErr {
					os.Remove(tt.filename)
				}
			}()

			// Вызываем тестируемую функцию
			createFile(tt.filename)

			// Проверяем, что файл создан
			if !tt.wantErr {
				// Проверяем существование файла
				if _, err := os.Stat(tt.filename); os.IsNotExist(err) {
					t.Errorf("createFile() failed to create file %s", tt.filename)
				}

				// Проверяем размер файла
				fileInfo, err := os.Stat(tt.filename)
				if err != nil {
					t.Errorf("failed to get file info: %v", err)
				}

				if fileInfo.Size() != PAGE_SIZE {
					t.Errorf("createFile() created file with wrong size, got %d, want %d", fileInfo.Size(), PAGE_SIZE)
				}

				// Читаем содержимое файла
				data, err := os.ReadFile(tt.filename)
				if err != nil {
					t.Errorf("failed to read file: %v", err)
				}

				// Десериализуем страницу
				page := deserializePage(data)

				// Проверяем корректность заголовка
				if page.Header.PageId != 4 {
					t.Errorf("wrong page id, got %d, want %d", page.Header.PageId, 1)
				}
				if page.Header.PageSize != PAGE_SIZE {
					t.Errorf("wrong page size, got %d, want %d", page.Header.PageSize, PAGE_SIZE)
				}

				// Проверяем количество слотов
				if len(page.Slots) != MAX_SLOTS {
					t.Errorf("wrong number of slots, got %d, want %d", len(page.Slots), MAX_SLOTS)
				}

				// Проверяем размер данных
				if len(page.Data) != DATA_SIZE {
					t.Errorf("wrong data size, got %d, want %d", len(page.Data), DATA_SIZE)
				}
			} else {
				// Для негативного теста проверяем, что файл НЕ создан
				if _, err := os.Stat(tt.filename); !os.IsNotExist(err) {
					t.Errorf("createFile() should not create file in non-existent directory")
				}
			}
		})
	}
}

func TestCreateFileInExistingDirectory(t *testing.T) {
	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "dbtest")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Создаем файл в существующей директории
	filename := filepath.Join(tmpDir, "test.db")
	createFile(filename)

	// Проверяем, что файл создан
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("createFile() failed to create file in existing directory")
	}
}

func TestCreateFileOverwrite(t *testing.T) {
	filename := "test_overwrite.db"

	// Создаем файл первый раз
	createFile(filename)
	defer os.Remove(filename)

	// Получаем информацию о первом файле
	firstFileInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("failed to get first file info: %v", err)
	}

	// Создаем файл второй раз (перезапись)
	createFile(filename)

	// Получаем информацию о втором файле
	secondFileInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("failed to get second file info: %v", err)
	}

	// Проверяем, что размер не изменился
	if firstFileInfo.Size() != secondFileInfo.Size() {
		t.Errorf("file sizes different after overwrite: first %d, second %d",
			firstFileInfo.Size(), secondFileInfo.Size())
	}
}
