package data

type PageHeader struct {
	PageId      uint32
	FreeSpace   uint16
	SlotsAmount uint16
}

type PageSlot struct {
	SlotId    uint16
	RowSize   uint16
	Offset    uint16
	IsDeleted bool
}

type Page struct {
	Header PageHeader
	Slots  []PageSlot
	Data   []DataRow
}

type DataCell struct {
	Value  interface{}
	Type   ColumnType
	IsNull bool
}

type DataRow struct {
	PageId uint32
	SlotId uint16
	Row    []DataCell
}

// Meta data

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

type MetaData struct {
	Name      string
	PageCount uint32
	Columns   []Column
}
