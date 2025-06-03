package disk_manager

import (
	"custom-database/config"
	"custom-database/internal/disk_manager/data"
	"custom-database/internal/disk_manager/meta"
	"fmt"
	"os"
	"path/filepath"
)

type InsertRowResult struct {
	SlotID uint16
	PageID uint32
}

type DiskManagerService interface {
	CreateTable(tableName string, columns []meta.Column) error
	InsertDataRow(tableName string, row []data.DataCell) (*InsertRowResult, error)
	ReadDataRow(tableName string, slotID uint16, pageID uint32) (*data.DataRow, error)
	DeleteDataRow(tableName string, slotID uint16, pageID uint32) error
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

func (dm *diskManager) InsertDataRow(tableName string, row []data.DataCell) (*InsertRowResult, error) {
	if err := dm.checkTableExists(tableName); err != nil {
		return nil, fmt.Errorf("InsertDataRow(): checkTableExists: %w", err)
	}

	filePath := filepath.Join(dm.cfg.DBPath, tableName)
	metaFile, err := meta.LoadMetaFile(tableName, filePath)
	if err != nil {
		return nil, fmt.Errorf("InsertDataRow(): meta.LoadMetaFile: %w", err)
	}

	fc, err := data.NewFileConnection(metaFile, filePath, false)
	if err != nil {
		return nil, fmt.Errorf("CreateTable(): data.NewFileConnection: %w", err)
	}
	defer fc.Close()

	result, err := fc.InsertDataRow(row)
	if err != nil {
		return nil, fmt.Errorf("InsertDataRow(): fc.InsertDataRow: %w", err)
	}

	return &InsertRowResult{SlotID: result.SlotID, PageID: result.PageID}, nil
}

func (dm *diskManager) ReadDataRow(tableName string, slotID uint16, pageID uint32) (*data.DataRow, error) {
	return nil, nil
}

func (dm *diskManager) DeleteDataRow(tableName string, slotID uint16, pageID uint32) error {
	return nil
}

func (dm *diskManager) checkTableExists(tableName string) error {
	filePath := filepath.Join(dm.cfg.DBPath, tableName)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("checkTableExists(): table not found: %w", err)
	}

	return nil
}
