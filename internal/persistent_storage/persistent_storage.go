package persistent_storage

import (
	"custom-database/config"
	"custom-database/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PersistentStorageService interface {
	CreateTable(tableName string, columns []models.Column) error
	DropTable(tableName string) error
	Insert(tableName string, newValues []models.Cell) error
	Select(tableName string) (*models.Table, error)
}

type persistentStorage struct {
	dir string
}

func NewPersistentStorage(cfg *config.Config) (PersistentStorageService, error) {
	if err := os.MkdirAll(cfg.DBPath, 0755); err != nil {
		return nil, fmt.Errorf("NewPersistentStorage(): failed to create persistent storage: %w", err)
	}

	return &persistentStorage{
		dir: cfg.DBPath,
	}, nil
}

func (ps *persistentStorage) CreateTable(tableName string, columns []models.Column) error {
	filename := filepath.Join(ps.dir, tableName+".json")

	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("CreateTable(): table already exists: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("CreateTable(): failed to create table file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(models.Table{
		Name:    tableName,
		Columns: columns,
	}); err != nil {
		return fmt.Errorf("CreateTable(): failed to encode table: %w", err)
	}

	return nil
}

func (ps *persistentStorage) DropTable(tableName string) error {
	filename := filepath.Join(ps.dir, tableName+".json")

	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("DropTable(): table does not exist: %w", err)
	}

	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("DropTable(): failed to remove table file: %w", err)
	}

	return nil
}

func (ps *persistentStorage) Insert(tableName string, allValues []models.Cell) error {
	filename := filepath.Join(ps.dir, tableName+".json")

	if _, err := os.Stat(filename); err != nil {
		return fmt.Errorf("Insert(): table does not exist: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("Insert(): failed to open table file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(allValues); err != nil {
		return fmt.Errorf("Insert(): failed to encode new values: %w", err)
	}

	return nil
}

func (ps *persistentStorage) Select(tableName string) (*models.Table, error) {
	filename := filepath.Join(ps.dir, tableName+".json")

	if _, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("Select(): table does not exist: %w", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Select(): failed to open table file: %w", err)
	}
	defer file.Close()

	var table models.Table
	if err := json.NewDecoder(file).Decode(&table); err != nil {
		return nil, fmt.Errorf("Select(): failed to decode table: %w", err)
	}

	return &table, nil
}
