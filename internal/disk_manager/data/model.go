package data

type PageHeader struct {
	PageId   uint32
	PageSize uint16
}

type PageSlot struct {
	RowId     uint16
	Offset    uint16
	Size      uint16
	IsDeleted bool
}

type Page struct {
	Header PageHeader
	Slots  []PageSlot
	Data   []byte
}
