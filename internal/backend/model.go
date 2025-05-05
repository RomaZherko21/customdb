package backend

import (
	"bytes"
	"custom-database/internal/models"
	"encoding/binary"
)

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

func (mc MemoryCell) AsBoolean() bool {
	if string(mc) == "true" {
		return true
	}

	return false
}

type table struct {
	columns     []string
	columnTypes []models.ColumnType
	rows        [][]MemoryCell
}
