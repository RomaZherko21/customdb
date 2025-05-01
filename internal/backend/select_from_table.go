package backend

import (
	"bytes"
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"encoding/binary"
	"fmt"
	"slices"
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

	//// FILTERING
	selectedColumnNames := []string{}
	for _, value := range statement.SelectedColumns {
		selectedColumnNames = append(selectedColumnNames, value.Literal.Value)
	}

	for i, row := range table.Rows {
		resultRow := []interface{}{}
		for i, column := range table.Columns {
			if slices.Contains(selectedColumnNames, column.Name) {
				resultRow = append(resultRow, row[i])
			}
		}

		table.Rows[i] = resultRow
	}
	newColumns := []models.Column{}
	for _, column := range table.Columns {
		if slices.Contains(selectedColumnNames, column.Name) {
			newColumns = append(newColumns, column)
		}
	}
	table.Columns = newColumns
	////

	rows := [][]models.Cell{}

	for _, row := range table.Rows {
		newRow := []models.Cell{}
		for i, cell := range row {
			column := newColumns[i]
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

			newRow = append(newRow, memoryCell)
		}
		rows = append(rows, newRow)
	}

	return &models.Table{
		Columns: newColumns,
		Rows:    rows,
	}, nil
}
