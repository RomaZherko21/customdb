package backend

import (
	"custom-database/config"
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"custom-database/internal/storage/memory_storage"
	"custom-database/internal/storage/persistent_storage"
)

type MemoryBackendService interface {
	ExecuteStatement(*ast.Ast) (*models.Table, error)
}

type memoryBackend struct {
	tables map[string]*table

	memoryStorage     memory_storage.MemoryStorageService
	persistentStorage persistent_storage.PersistentStorageService
}

func NewMemoryBackend(config *config.Config) (MemoryBackendService, error) {
	memoryStorage := memory_storage.NewMemoryStorage()
	persistentStorage, err := persistent_storage.NewPersistentStorage(config)
	if err != nil {
		return nil, err
	}

	return &memoryBackend{
		memoryStorage:     memoryStorage,
		persistentStorage: persistentStorage,
		tables:            map[string]*table{},
	}, nil
}

func (mb *memoryBackend) ExecuteStatement(a *ast.Ast) (*models.Table, error) {
	var err error

	for _, stmt := range a.Statements {
		switch stmt.Kind {
		case ast.CreateTableKind:
			err = mb.createTable(stmt.CreateTableStatement)
			if err != nil {
				return nil, err
			}
		case ast.DropTableKind:
			err = mb.dropTable(stmt.DropTableStatement)
			if err != nil {
				return nil, err
			}
		case ast.InsertKind:
			err = mb.insertIntoTable(stmt.InsertStatement)
			if err != nil {
				return nil, err
			}
		case ast.SelectKind:
			results, err := mb.selectFromTable(stmt.SelectStatement)
			if err != nil {
				return nil, err
			}

			return results, nil
		}
	}

	return nil, nil
}
