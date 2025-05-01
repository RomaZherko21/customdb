package backend

import (
	"bytes"
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"encoding/binary"
	"fmt"
)

func (mb *memoryBackend) selectFromTable(statement *ast.SelectStatement) (*models.Table, error) {
	// table, err := mb.memoryStorage.Select(statement.From.Value)
	// if err != nil {
	// 	return nil, err
	// }

	table, err := mb.persistentStorage.Select(statement.From.Value)
	if err != nil {
		return nil, err
	}

	results := [][]models.Cell{}

	for _, row := range table.Rows {
		result := []models.Cell{}
		for i, cell := range row {
			column := table.Columns[i]
			var memoryCell MemoryCell

			if column.Type == models.IntType {
				buf := new(bytes.Buffer)
				err := binary.Write(buf, binary.BigEndian, int32(cell.(float64)))
				if err != nil {
					return nil, fmt.Errorf("failed to convert int: %w", err)
				}
				memoryCell = MemoryCell(buf.Bytes())
			}

			if column.Type == models.TextType {
				memoryCell = MemoryCell(cell.(string))
			}

			result = append(result, memoryCell)
		}
		results = append(results, result)
	}

	return &models.Table{
		Columns: table.Columns,
		Rows:    results,
	}, nil
}
