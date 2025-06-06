package data

import (
	"fmt"
	"os"
)

type Service interface {
	InsertDataRow(row []DataCell) (*InsertRowResult, error)
}

type dataService struct {
	file          *os.File
	tableName     string
	meta          *MetaData
	metaDataSpace int
}

// NewDataService на вход получает файл, который будет использоваться для работы с данными
func NewDataService(file *os.File, tableName string) (Service, error) {
	if file == nil {
		return nil, fmt.Errorf("NewDataService(): file is nil")
	}

	ds := &dataService{
		file:      file,
		tableName: tableName,
	}

	metaData, metaDataSpace, err := ds.loadMetaData()
	if err != nil {
		return nil, fmt.Errorf("NewDataService(): loadMetaData: %w", err)
	}
	ds.meta = metaData
	ds.metaDataSpace = metaDataSpace

	return ds, nil
}

// InitTableData создает новую таблицу в файле
func InitTableData(file *os.File, tableName string, columns []Column) (Service, error) {
	ds := &dataService{
		file:      file,
		tableName: tableName,
	}

	metaData, metaDataSpace, err := ds.writeMetaData(ds.tableName, columns)
	if err != nil {
		return nil, fmt.Errorf("InitTableData(): WriteMetaData: %w", err)
	}
	ds.meta = metaData
	ds.metaDataSpace = metaDataSpace

	err = ds.writeInitialPageData()
	if err != nil {
		return nil, fmt.Errorf("InitTableData(): WriteInitialPageData: %w", err)
	}

	return ds, nil
}

type InsertRowResult struct {
	SlotID uint16
	PageID uint32
}

func (ds *dataService) InsertDataRow(row []DataCell) (*InsertRowResult, error) {
	return ds.insertDataRow(row)
}
