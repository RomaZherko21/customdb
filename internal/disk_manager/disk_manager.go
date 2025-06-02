package disk_manager

import (
	"custom-database/config"
	"custom-database/internal/disk_manager/data"
	"custom-database/internal/disk_manager/meta"
	"fmt"
	"os"
	"path/filepath"
)

type DiskManagerService interface {
	CreateTable(filename string, columns []meta.Column) error
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

func (dm *diskManager) CreateTable(tableName string, columns []meta.Column) error {
	filePath := filepath.Join(dm.cfg.DBPath, tableName)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("CreateTable(): table already exists: %w", err)
	}

	if err := os.MkdirAll(filePath, 0755); err != nil {
		return fmt.Errorf("CreateTable(): os.MkdirAll: %w", err)
	}

	metaFile, err := meta.CreateMetaFile(tableName, columns, filePath)
	if err != nil {
		return fmt.Errorf("CreateTable(): meta.CreateMetaFile: %w", err)
	}

	fc, err := data.NewFileConnection(metaFile, filePath, true)
	if err != nil {
		return fmt.Errorf("CreateTable(): data.NewFileConnection: %w", err)
	}
	defer fc.Close()

	return nil
}
