package data

import "custom-database/internal/disk_manager/meta"

type PageHeader struct {
	PageId    uint32
	FreeSpace uint16
}

type PageSlot struct {
	RowId     uint16
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
	Type   meta.ColumnType
	IsNull bool
}

type DataRow struct {
	PageId uint32
	SlotId uint16
	Row    []DataCell
}
