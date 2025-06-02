package meta

type ColumnType uint8

const (
	TypeInt32 ColumnType = iota
	TypeInt64
	TypeUint32
	TypeUint64
	TypeBoolean
	TypeText
)

type Column struct {
	Name       string
	Type       ColumnType
	IsNullable bool
}

type MetaFile struct {
	Name      string
	PageCount uint32
	Columns   []Column
}
