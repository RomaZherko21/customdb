package disk_manager

import (
	"custom-database/config"
	"custom-database/internal/disk_manager/data"
	"fmt"
	"os"
	"path/filepath"
)

type DiskManagerService interface {
	CreateTable(tableName string, columns []data.Column) error
}

type diskManager struct {
	cfg *config.Config
}

func NewDiskManager(cfg *config.Config) (DiskManagerService, error) {
	if err := os.MkdirAll(cfg.DBPath, 0755); err != nil {
		return nil, fmt.Errorf("NewDiskManager(): failed to create folder in disk_manager: %v", err)
	}

	return &diskManager{
		cfg: cfg,
	}, nil
}

func (dm *diskManager) CreateTable(tableName string, columns []data.Column) error {
	filePath := filepath.Join(dm.cfg.DBPath, tableName)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("CreateTable(): table already exists: %w", err)
	}

	if err := os.MkdirAll(filePath, 0755); err != nil {
		return fmt.Errorf("CreateTable(): os.MkdirAll: %w", err)
	}

	file, err := os.Create(filepath.Join(filePath, tableName+".data"))
	if err != nil {
		return fmt.Errorf("CreateTable(): os.Create: %w", err)
	}

	_, err = data.InitTableData(file, tableName, columns)
	if err != nil {
		return fmt.Errorf("CreateTable(): data.InitTableData: %w", err)
	}

	file.Close()

	return nil
}
