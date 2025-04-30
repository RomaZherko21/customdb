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

type MemoryBackendService interface {
	ExecuteStatement(*ast.Ast) (*models.Table, error)
}

type memoryBackend struct {
	tables map[string]*table
}

func NewMemoryBackend() MemoryBackendService {
	return &memoryBackend{
		tables: map[string]*table{},
	}
}

func (mb *memoryBackend) ExecuteStatement(a *ast.Ast) (*models.Table, error) {
	var err error

	for _, stmt := range a.Statements {
		switch stmt.Kind {
		case ast.CreateTableKind:
			err = mb.createTable(a.Statements[0].CreateTableStatement)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		case ast.DropTableKind:
			err = mb.dropTable(stmt.DropTableStatement)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		case ast.InsertKind:
			err = mb.Insert(stmt.InsertStatement)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		case ast.SelectKind:
			results, err := mb.Select(stmt.SelectStatement)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			return results, nil
		}
	}

	return nil, nil
}

func (mb *memoryBackend) createTable(statement *ast.CreateTableStatement) error {
	t := table{}
	mb.tables[statement.Name.Value] = &t
	if statement.Cols == nil {
		return nil
	}

	for _, col := range *statement.Cols {
		t.columns = append(t.columns, col.Name.Value)

		var dt models.ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = models.IntType
		case "text":
			dt = models.TextType
		default:
			return fmt.Errorf("Invalid datatype: %s", col.Datatype.Value)
		}

		t.columnTypes = append(t.columnTypes, dt)
	}

	return nil
}

func (mb *memoryBackend) dropTable(statement *ast.DropTableStatement) error {
	_, ok := mb.tables[statement.Table.Value]
	if !ok {
		return fmt.Errorf("Table does not exist: %s", statement.Table.Value)
	}

	delete(mb.tables, statement.Table.Value)
	return nil
}

func (mb *memoryBackend) Insert(statement *ast.InsertStatement) error {
	table, ok := mb.tables[statement.Table.Value]
	if !ok {
		return fmt.Errorf("Table does not exist: %s", statement.Table.Value)
	}

	if statement.Values == nil {
		return nil
	}

	row := []MemoryCell{}

	if len(*statement.Values) != len(table.columns) {
		return fmt.Errorf("Missing values: %d != %d", len(*statement.Values), len(table.columns))
	}

	for _, value := range *statement.Values {
		if value.Kind != ast.LiteralKind {
			fmt.Println("Skipping non-literal.")
			continue
		}

		row = append(row, mb.tokenToCell(value.Literal))
	}

	table.rows = append(table.rows, row)
	return nil
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

func (mb *memoryBackend) Select(statement *ast.SelectStatement) (*models.Table, error) {
	table, ok := mb.tables[statement.From.Value]
	if !ok {
		return nil, fmt.Errorf("Table does not exist: %s", statement.From.Value)
	}

	results := [][]models.Cell{}
	columns := []models.Column{}

	for i, row := range table.rows {
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
				for i, tableCol := range table.columns {
					if tableCol == lit.Value {
						if isFirstRow {
							columns = append(columns, models.Column{
								Type: table.columnTypes[i],
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
