package model

type Column struct {
	Name string
	Type DataType
}

type Row struct {
	Values []interface{}
}

type Table struct {
	TableName string
	Columns   []Column
	Rows      []Row
}
