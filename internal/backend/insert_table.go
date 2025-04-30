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
