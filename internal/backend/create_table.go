package backend

import (
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
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
		case "int":
			dt = models.IntType
		case "text":
			dt = models.TextType
		default:
			return fmt.Errorf("Invalid datatype: %s", col.Datatype.Value)
		}

		columns = append(columns, models.Column{
			Name: col.Name.Value,
			Type: dt,
		})
	}

	mb.memoryStorage.CreateTable(statement.Name.Value, columns)

	return nil
}
