package backend

import (
	"bytes"
	"custom-database/internal/parser/ast"
	"custom-database/internal/parser/lex"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

type ColumnType uint

type Column struct {
	Name string
	Type ColumnType
}

const (
	TextType ColumnType = iota
	IntType
)

type Cell interface {
	AsText() string
	AsInt() int32
}

type Results struct {
	Columns []Column
	Rows    [][]Cell
}

var (
	ErrTableDoesNotExist  = errors.New("Table does not exist")
	ErrColumnDoesNotExist = errors.New("Column does not exist")
	ErrInvalidSelectItem  = errors.New("Select item is not valid")
	ErrInvalidDatatype    = errors.New("Invalid datatype")
	ErrMissingValues      = errors.New("Missing values")
)

type Backend interface {
	CreateTable(*ast.CreateTableStatement) error
	Insert(*ast.InsertStatement) error
	Select(*ast.SelectStatement) (*Results, error)
}

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}

	return i
}

func (mc MemoryCell) AsText() string {
	return string(mc)
}

type table struct {
	columns     []string
	columnTypes []ColumnType
	rows        [][]MemoryCell
}

type MemoryBackend struct {
	tables map[string]*table
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		tables: map[string]*table{},
	}
}

func (mb *MemoryBackend) CreateTable(statement *ast.CreateTableStatement) error {
	t := table{}
	mb.tables[statement.Name.Value] = &t
	if statement.Cols == nil {
		return nil
	}

	for _, col := range *statement.Cols {
		t.columns = append(t.columns, col.Name.Value)

		var dt ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		default:
			return ErrInvalidDatatype
		}

		t.columnTypes = append(t.columnTypes, dt)
	}

	return nil
}

func (mb *MemoryBackend) DropTable(statement *ast.DropTableStatement) error {
	_, ok := mb.tables[statement.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}

	delete(mb.tables, statement.Table.Value)
	return nil
}

func (mb *MemoryBackend) Insert(statement *ast.InsertStatement) error {
	table, ok := mb.tables[statement.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}

	if statement.Values == nil {
		return nil
	}

	row := []MemoryCell{}

	if len(*statement.Values) != len(table.columns) {
		return ErrMissingValues
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

func (mb *MemoryBackend) tokenToCell(t *lex.Token) MemoryCell {
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

func (mb *MemoryBackend) Select(statement *ast.SelectStatement) (*Results, error) {
	table, ok := mb.tables[statement.From.Value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}

	results := [][]Cell{}
	columns := []Column{}

	for i, row := range table.rows {
		result := []Cell{}
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
							columns = append(columns, Column{
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
					return nil, ErrColumnDoesNotExist
				}

				continue
			}

			return nil, ErrColumnDoesNotExist
		}

		results = append(results, result)
	}

	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}
