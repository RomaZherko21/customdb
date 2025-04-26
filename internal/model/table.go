package model

type Column struct {
	Name string
	Type DataType
}

type Table struct {
	TableName string
	Columns   []Column
}
