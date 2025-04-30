package models

import (
	"bytes"
	"encoding/binary"
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

type Table struct {
	Name    string
	Columns []Column
	Rows    [][]Cell
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
