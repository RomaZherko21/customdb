package data

type PageHeader struct {
	PageId   int32
	PageSize int32
}

type PageSlot struct {
	RowId     int32
	Offset    int32
	Size      int32
	IsDeleted bool
}

type Page struct {
	Header *PageHeader
	Slots  []PageSlot
	Data   []byte
}
