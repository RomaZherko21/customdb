package models

type ColumnType uint

type Column struct {
	Name string     `json:"name"`
	Type ColumnType `json:"type"`
}

const (
	TextType ColumnType = iota
	IntType
	BoolType
	TimestampType
)

type Cell interface {
	AsText() string
	AsInt() int32
	AsBoolean() bool
	IsNull() bool
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}
