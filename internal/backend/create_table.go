package backend

import (
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"fmt"
)

func (mb *memoryBackend) createTable(statement *ast.CreateTableStatement) error {
	if statement.Cols == nil {
		return nil
	}

	columns := []models.Column{}
	for _, col := range *statement.Cols {
		var dt models.ColumnType

		switch col.Datatype.Value {
		case string(lex.IntKeyword):
			dt = models.IntType
		case string(lex.TextKeyword):
			dt = models.TextType
		case string(lex.BooleanTypeKeyword):
			dt = models.BoolType
		case string(lex.TimestampTypeKeyword):
			dt = models.TimestampType
		default:
			return fmt.Errorf("Invalid datatype: %s", col.Datatype.Value)
		}

		columns = append(columns, models.Column{
			Name: col.Name.Value,
			Type: dt,
		})
	}

	err := mb.persistentStorage.CreateTable(statement.Name.Value, columns)
	if err != nil {
		return err
	}

	// return mb.memoryStorage.CreateTable(statement.Name.Value, columns)
	return nil
}
