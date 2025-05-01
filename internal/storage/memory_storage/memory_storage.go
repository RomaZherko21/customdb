package memory_storage

import (
	"custom-database/internal/models"
	"fmt"
)

type MemoryStorageService interface {
	CreateTable(tableName string, columns []models.Column) error
	DropTable(tableName string) error
	Insert(tableName string, newValues []models.Cell) error
	Select(tableName string) (*models.Table, error)
}

type memoryStorage struct {
	tables map[string]*models.Table
}

func NewMemoryStorage() MemoryStorageService {
	return &memoryStorage{
		tables: map[string]*models.Table{},
	}
}

func (ms *memoryStorage) CreateTable(tableName string, columns []models.Column) error {
	if _, ok := ms.tables[tableName]; ok {
		return fmt.Errorf("CreateTable(): table already exists")
	}

	ms.tables[tableName] = &models.Table{
		Name:    tableName,
		Columns: columns,
	}

	return nil
}

func (ms *memoryStorage) DropTable(tableName string) error {
	_, ok := ms.tables[tableName]
	if !ok {
		return fmt.Errorf("DropTable(): table does not exist")
	}

	delete(ms.tables, tableName)
	return nil
}

func (ms *memoryStorage) Insert(tableName string, newValues []models.Cell) error {
	table, ok := ms.tables[tableName]
	if !ok {
		return fmt.Errorf("Insert(): table does not exist")
	}

	if len(newValues) != len(table.Columns) {
		return fmt.Errorf("Insert(): missing values")
	}

	table.Rows = append(table.Rows, newValues)
	return nil
}

func (ms *memoryStorage) Select(tableName string) (*models.Table, error) {
	table, ok := ms.tables[tableName]
	if !ok {
		return nil, fmt.Errorf("Select(): table does not exist")
	}

	return table, nil
}
