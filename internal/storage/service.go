package storage

import (
	"custom-database/internal/model"
	"fmt"
)

type Storage interface {
	GetTable(name string) storageTable
	CreateTable(table model.Table)
	InsertInto(table model.Table)
}

type storageTable struct {
	data    []interface{}
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

func (s *storage) CreateTable(table model.Table) {
	s.tables[table.TableName] = storageTable{
		data:    nil,
		columns: table.Columns,
	}
}

func (s *storage) InsertInto(table model.Table) {
	fmt.Println(table)
}
