package backend

import (
	"bytes"
	"custom-database/internal/models"
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"encoding/binary"
	"fmt"
	"strconv"
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

func (mb *memoryBackend) dropTable(statement *ast.DropTableStatement) error {
	return mb.memoryStorage.DropTable(statement.Table.Value)
}

func (mb *memoryBackend) insertIntoTable(statement *ast.InsertStatement) error {
	if statement.Values == nil {
		return nil
	}

	row := []MemoryCell{}
	for _, value := range *statement.Values {
		if value.Kind != ast.LiteralKind {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, mb.tokenToCell(value.Literal))
	}

	cells := make([]models.Cell, len(row))
	for i, cell := range row {
		cells[i] = cell
	}

	return mb.memoryStorage.Insert(statement.Table.Value, cells)
}

func (mb *memoryBackend) tokenToCell(t *lex.Token) MemoryCell {
	if t.Kind == lex.NumericToken {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			panic(err)
		}

		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return MemoryCell(buf.Bytes())
	}

	if t.Kind == lex.StringToken {
		return MemoryCell(t.Value)
	}

	return nil
}

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
