package executor

import (
	"custom-database/internal/model"
	"custom-database/internal/storage"
)

type Executor interface {
	CreateTable(command model.Table) error
	InsertInto(command model.Table) error
}

type executor struct {
	storage storage.Storage
}

func NewExecutor(storage storage.Storage) Executor {
	return &executor{
		storage: storage,
	}
}

func (e *executor) CreateTable(command model.Table) error {
	return e.storage.CreateTable(command)
}

func (e *executor) InsertInto(command model.Table) error {
	return e.storage.InsertInto(command)
}
