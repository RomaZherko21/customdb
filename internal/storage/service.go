package storage

import (
	"custom-database/internal/model"
	"fmt"
)

type Storage interface {
	GetTable(name string) storageTable
	CreateTable(table model.Table) error
	InsertInto(table model.Table) error
}

type storageTable struct {
	rows    [][]interface{}
	columns []model.Column
}

type storage struct {
	tables map[string]storageTable
}

func NewStorage() Storage {
	return &storage{
		tables: make(map[string]storageTable),
	}
}

func (s *storage) GetTable(name string) storageTable {
	return s.tables[name]
}

func (s *storage) CreateTable(table model.Table) error {
	s.tables[table.TableName] = storageTable{
		rows:    [][]interface{}{},
		columns: table.Columns,
	}

	return nil
}

func (s *storage) InsertInto(table model.Table) error {

	tableName, ok := s.tables[table.TableName]
	if !ok {
		return fmt.Errorf("table %s not found", table.TableName)
	}

	tableName.rows = append(tableName.rows, table.Rows[0])

	s.tables[table.TableName] = tableName

	return nil
}
