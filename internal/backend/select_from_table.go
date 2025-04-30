package backend

import (
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"fmt"
)

func (mb *memoryBackend) selectFromTable(statement *ast.SelectStatement) (*models.Table, error) {
	table, err := mb.memoryStorage.Select(statement.From.Value)
	if err != nil {
		return nil, err
	}

	results := [][]models.Cell{}
	columns := []models.Column{}

	for i, row := range table.Rows {
		result := []models.Cell{}
		isFirstRow := i == 0

		for _, exp := range statement.Item {
			if exp.Kind != ast.LiteralKind {
				// Unsupported, doesn't currently exist, ignore.
				fmt.Println("Skipping non-literal expression.")
				continue
			}

			lit := exp.Literal
			if lit.Kind == lex.IdentifierToken {
				found := false
				for i, tableCol := range table.Columns {
					if tableCol.Name == lit.Value {
						if isFirstRow {
							columns = append(columns, models.Column{
								Type: tableCol.Type,
								Name: lit.Value,
							})
						}

						result = append(result, row[i])
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("Column does not exist: %s", lit.Value)
				}

				continue
			}

			return nil, fmt.Errorf("Column does not exist: %s", lit.Value)
		}

		results = append(results, result)
	}

	return &models.Table{
		Columns: columns,
		Rows:    results,
	}, nil
}
