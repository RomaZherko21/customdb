package meta

import (
	"fmt"
	"os"
	"testing"

	bs "custom-database/internal/disk_manager/binary_serializer"
)

func TestCreateMetaFile(t *testing.T) {
	t.Run("Create Simple Meta File", func(t *testing.T) {
		// Подготовка тестовых данных
		metaFile := &MetaFile{
			Name: "users",
			Columns: []Column{
				{Name: "id", Type: TypeInt32, IsNullable: false},
				{Name: "name", Type: TypeText, IsNullable: true},
				{Name: "surname", Type: TypeText, IsNullable: true},
				{Name: "age", Type: TypeUint32, IsNullable: true},
				{Name: "height", Type: TypeUint64, IsNullable: true},
				{Name: "residence", Type: TypeText, IsNullable: false},
				{Name: "is_admin", Type: TypeBoolean, IsNullable: false},
			},
		}

		// Создаем файл
		err := CreateMetaFile(metaFile, "")
		if err != nil {
			t.Fatalf("Failed to create meta file: %v", err)
		}

		// Проверяем что файл создан
		fileName := metaFile.Name + ".meta"
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Errorf("Meta file %s was not created", fileName)
		}

		// Читаем содержимое файла
		data, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("Failed to read meta file: %v", err)
		}

		// Десериализуем и проверяем содержимое
		deserializedMeta := deserializeMetaFile(data)

		// Проверяем имя таблицы
		if deserializedMeta.Name != metaFile.Name {
			t.Errorf("Table name mismatch: got %s, want %s",
				deserializedMeta.Name, metaFile.Name)
		}

		// Проверяем количество колонок
		if len(deserializedMeta.Columns) != len(metaFile.Columns) {
			t.Errorf("Columns count mismatch: got %d, want %d",
				len(deserializedMeta.Columns), len(metaFile.Columns))
		}

		// Проверяем каждую колонку
		for i, col := range metaFile.Columns {
			deserializedCol := deserializedMeta.Columns[i]
			if deserializedCol.Name != col.Name {
				t.Errorf("Column %d name mismatch: got %s, want %s",
					i, deserializedCol.Name, col.Name)
			}
			if deserializedCol.Type != col.Type {
				t.Errorf("Column %d type mismatch: got %d, want %d",
					i, deserializedCol.Type, col.Type)
			}
		}

		// Очистка: удаляем созданный файл
		os.Remove(fileName)
	})

	t.Run("Meta File with Empty Columns", func(t *testing.T) {
		metaFile := &MetaFile{
			Name:    "empty_table",
			Columns: []Column{},
		}

		err := CreateMetaFile(metaFile, "")
		if err != nil {
			t.Fatalf("Failed to create meta file: %v", err)
		}

		fileName := metaFile.Name + ".meta"
		data, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("Failed to read meta file: %v", err)
		}

		deserializedMeta := deserializeMetaFile(data)

		if deserializedMeta.Name != metaFile.Name {
			t.Errorf("Table name mismatch: got %s, want %s",
				deserializedMeta.Name, metaFile.Name)
		}

		if len(deserializedMeta.Columns) != 0 {
			t.Errorf("Expected empty columns, got %d columns",
				len(deserializedMeta.Columns))
		}

		os.Remove(fileName)
	})

	t.Run("Meta File with Long Names", func(t *testing.T) {
		metaFile := &MetaFile{
			Name: "very_long_table_name_that_tests_string_serialization",
			Columns: []Column{
				{Name: "very_long_column_name_to_test_string_serialization_and_deserialization", Type: TypeInt32},
				{Name: "another_very_long_column_name_to_ensure_proper_handling_of_long_strings", Type: TypeText},
			},
		}

		err := CreateMetaFile(metaFile, "")
		if err != nil {
			t.Fatalf("Failed to create meta file: %v", err)
		}

		fileName := metaFile.Name + ".meta"
		data, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("Failed to read meta file: %v", err)
		}

		deserializedMeta := deserializeMetaFile(data)

		if deserializedMeta.Name != metaFile.Name {
			t.Errorf("Table name mismatch: got %s, want %s",
				deserializedMeta.Name, metaFile.Name)
		}

		for i, col := range metaFile.Columns {
			deserializedCol := deserializedMeta.Columns[i]
			if deserializedCol.Name != col.Name {
				t.Errorf("Column %d name mismatch: got %s, want %s",
					i, deserializedCol.Name, col.Name)
			}
		}

		os.Remove(fileName)
	})

	t.Run("Meta File with All Data Types", func(t *testing.T) {
		metaFile := &MetaFile{
			Name: "all_types",
			Columns: []Column{
				{Name: "int32_col", Type: TypeInt32},
				{Name: "int64_col", Type: TypeInt64},
				{Name: "uint32_col", Type: TypeUint32},
				{Name: "uint64_col", Type: TypeUint64},
				{Name: "bool_col", Type: TypeBoolean},
				{Name: "string_col", Type: TypeText},
			},
		}

		err := CreateMetaFile(metaFile, "")
		if err != nil {
			t.Fatalf("Failed to create meta file: %v", err)
		}

		fileName := metaFile.Name + ".meta"
		data, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("Failed to read meta file: %v", err)
		}

		deserializedMeta := deserializeMetaFile(data)

		if len(deserializedMeta.Columns) != len(metaFile.Columns) {
			t.Errorf("Columns count mismatch: got %d, want %d",
				len(deserializedMeta.Columns), len(metaFile.Columns))
		}

		for i, col := range metaFile.Columns {
			deserializedCol := deserializedMeta.Columns[i]
			if deserializedCol.Type != col.Type {
				t.Errorf("Column %d type mismatch: got %d, want %d",
					i, deserializedCol.Type, col.Type)
			}
		}

		os.Remove(fileName)
	})
}

func TestNullableColumns(t *testing.T) {
	t.Run("Serialize and Deserialize Nullable Columns", func(t *testing.T) {
		metaFile := &MetaFile{
			Name: "users",
			Columns: []Column{
				{Name: "id", Type: TypeInt32, IsNullable: false},
				{Name: "name", Type: TypeText, IsNullable: true},
				{Name: "email", Type: TypeText, IsNullable: true},
				{Name: "age", Type: TypeInt32, IsNullable: true},
				{Name: "is_admin", Type: TypeBoolean, IsNullable: false},
				{Name: "is_premium", Type: TypeBoolean, IsNullable: false},
				{Name: "is_deleted", Type: TypeBoolean, IsNullable: true},
			},
		}

		// Сериализуем
		data := serializeMetaFile(metaFile)

		offset := bs.TEXT_TYPE_HEADER + len(metaFile.Name) + COLUMN_COUNT_SIZE
		nullBitmap := bs.ReadUint32(data, offset)

		// Проверяем, что нужные биты установлены
		expectedNullable := []int{1, 2, 3, 6} // индексы nullable колонок
		for i := 0; i < 8; i++ {
			isNullable := false
			for _, idx := range expectedNullable {
				if i == idx {
					isNullable = true
					break
				}
			}
			if getBit(nullBitmap, i) != isNullable {
				t.Errorf("Column %d nullable status: got %v, want %v",
					i, getBit(nullBitmap, i), isNullable)
			}
		}

		// Десериализуем и проверяем
		deserializedMeta := deserializeMetaFile(data)

		// Проверяем, что все колонки правильно помечены как nullable
		for i, col := range metaFile.Columns {
			deserializedCol := deserializedMeta.Columns[i]
			if deserializedCol.IsNullable != col.IsNullable {
				t.Errorf("Column %d (%s) nullable status: got %v, want %v",
					i, col.Name, deserializedCol.IsNullable, col.IsNullable)
			}
		}
	})

	t.Run("Empty Table with Nullable Flag", func(t *testing.T) {
		metaFile := &MetaFile{
			Name:    "empty_table",
			Columns: []Column{},
		}

		data := serializeMetaFile(metaFile)
		deserializedMeta := deserializeMetaFile(data)

		if len(deserializedMeta.Columns) != 0 {
			t.Errorf("Expected empty columns, got %d columns", len(deserializedMeta.Columns))
		}
	})

	t.Run("Max Columns with Mixed Nullable", func(t *testing.T) {
		// Создаем таблицу с максимальным количеством колонок
		columns := make([]Column, MAX_COLUMNS)
		for i := 0; i < MAX_COLUMNS; i++ {
			columns[i] = Column{
				Name:       fmt.Sprintf("col_%d", i),
				Type:       TypeInt32,
				IsNullable: i%2 == 0, // Четные колонки nullable
			}
		}

		metaFile := &MetaFile{
			Name:    "max_columns",
			Columns: columns,
		}

		data := serializeMetaFile(metaFile)
		deserializedMeta := deserializeMetaFile(data)

		for i, col := range metaFile.Columns {
			deserializedCol := deserializedMeta.Columns[i]
			if deserializedCol.IsNullable != col.IsNullable {
				t.Errorf("Column %d nullable status: got %v, want %v",
					i, deserializedCol.IsNullable, col.IsNullable)
			}
		}
	})
}
