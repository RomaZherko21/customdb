package storage

import (
	"custom-database/config"
	"custom-database/internal/model"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

type Storage interface {
	GetTable(name string) storageTable
	CreateTable(table model.Table) error
	InsertInto(table model.Table) error
	Select(table model.Table) ([][]interface{}, error)
}

type storageTable struct {
	Rows    [][]interface{}
	Columns []model.Column
}

type storage struct {
	tables map[string]storageTable
	dir    string
}

func NewStorage(cfg *config.Config) Storage {
	if err := os.MkdirAll(cfg.DBPath, 0755); err != nil {
		panic(fmt.Sprintf("failed to create tables directory: %v", err))
	}

	return &storage{
		tables: make(map[string]storageTable),
		dir:    cfg.DBPath,
	}
}

func (s *storage) GetTable(name string) storageTable {
	return s.tables[name]
}

func (s *storage) CreateTable(table model.Table) error {
	// Создаем таблицу в памяти
	s.tables[table.TableName] = storageTable{
		Rows:    [][]interface{}{},
		Columns: table.Columns,
	}

	// Сохраняем таблицу в бинарный файл
	filename := filepath.Join(s.dir, table.TableName+".bin")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create table file: %w", err)
	}
	defer file.Close()

	// Создаем encoder и сохраняем таблицу
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(s.tables[table.TableName]); err != nil {
		return fmt.Errorf("failed to encode table: %w", err)
	}

	return nil
}

func (s *storage) InsertInto(table model.Table) error {
	filename := filepath.Join(s.dir, table.TableName+".bin")

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open table file: %w", err)
	}
	defer file.Close()

	// Декодируем данные из файла
	var tableData storageTable
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&tableData); err != nil {
		return fmt.Errorf("failed to decode table data: %w", err)
	}

	s.tables[table.TableName] = tableData

	tableName, ok := s.tables[table.TableName]
	if !ok {
		return fmt.Errorf("table %s not found", table.TableName)
	}

	tableName.Rows = append(tableName.Rows, table.Rows[0])

	s.tables[table.TableName] = tableName

	// Сохраняем обновленную таблицу в файл
	filename = filepath.Join(s.dir, table.TableName+".bin")
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to update table file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(tableName); err != nil {
		return fmt.Errorf("failed to encode updated table: %w", err)
	}

	return nil
}

func (s *storage) Select(table model.Table) ([][]interface{}, error) {
	// Проверяем существование файла таблицы
	filename := filepath.Join(s.dir, table.TableName+".bin")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("table %s not found", table.TableName)
	}

	// Открываем файл для чтения
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open table file: %w", err)
	}
	defer file.Close()

	// Декодируем данные из файла
	var tableData storageTable
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&tableData); err != nil {
		return nil, fmt.Errorf("failed to decode table data: %w", err)
	}

	// Если запрошены все колонки (*)
	if len(table.Columns) == 0 {
		return tableData.Rows, nil
	}

	// Если запрошены конкретные колонки
	result := make([][]interface{}, 0)
	for _, row := range tableData.Rows {
		// Создаем новую строку только с запрошенными колонками
		newRow := make([]interface{}, len(table.Columns))
		for i, col := range table.Columns {
			// Ищем индекс запрошенной колонки в исходной таблице
			for j, origCol := range tableData.Columns {
				if origCol.Name == col.Name {
					newRow[i] = row[j]
					break
				}
			}
		}
		result = append(result, newRow)
	}

	return result, nil
}
