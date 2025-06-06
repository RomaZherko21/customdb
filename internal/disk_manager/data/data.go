package data

import (
	"fmt"
)

func (ds *dataService) insertDataRow(row []DataCell) (*InsertRowResult, error) {
	pageHeaders, err := ds.ParsePageHeaders()
	if err != nil {
		return nil, fmt.Errorf("insertDataRow(): ds.ParsePageHeaders: %w", err)
	}

	rowSize := CalculateDataRowSize(row)
	currentPageID := uint32(0)

	for _, pageHeader := range pageHeaders {
		if uint32(pageHeader.FreeSpace) >= rowSize {
			currentPageID = pageHeader.PageId
			break
		}
	}
	if currentPageID == 0 {
		// TODO: handle case when no free space is found, and we need a new page
	}

	serializedRow := serializeDataRow(row)

	slot, err := ds.insertPageSlot(currentPageID, uint16(rowSize))
	if err != nil {
		return nil, fmt.Errorf("Insert(): InsertPageSlot: %w", err)
	}

	writeDataPosition := slot.Offset
	_, err = ds.file.WriteAt(serializedRow, int64(writeDataPosition))
	if err != nil {
		return nil, fmt.Errorf("Insert(): file.WriteAt: %w", err)
	}

	return &InsertRowResult{SlotID: slot.SlotId, PageID: currentPageID}, nil
}

func (ds *dataService) insertPageSlot(pageID uint32, rowSize uint16) (*PageSlot, error) {
	pageStartingPosition := ds.CalculatePageStartingPosition(pageID)

	slots, err := ds.ParsePageSlots(pageID)
	if err != nil {
		return nil, fmt.Errorf("insertPageSlot(): ds.ParsePageSlots: %w", err)
	}

	lastSlotOffset := uint16(PAGE_SIZE)
	if len(slots) > 0 {
		lastSlotOffset = slots[len(slots)-1].Offset
	}

	newSlotID := len(slots) + 1

	slot := PageSlot{
		SlotId:    uint16(newSlotID),
		Offset:    lastSlotOffset - rowSize,
		RowSize:   rowSize,
		IsDeleted: false,
	}

	serializedSlot := serializePageSlots([]PageSlot{slot})

	writeSlotPosition := pageStartingPosition + PAGE_HEADER_SIZE + uint32(len(slots)*ONE_SLOT_SIZE)
	ds.file.WriteAt(serializedSlot, int64(writeSlotPosition))

	return &slot, nil
}

func (fc *dataService) writeInitialPageData() error {
	page := &Page{
		Header: PageHeader{
			PageId:      INITIAL_PAGE_ID,
			FreeSpace:   DATA_SPACE,
			SlotsAmount: 0,
		},
		Slots: make([]PageSlot, MAX_SLOTS),
		Data:  make([]DataRow, 0),
	}
	data := serializePage(page)

	// Записываем начальную страницу после метаданных
	_, err := fc.file.WriteAt(data, int64(fc.metaDataSpace))
	if err != nil {
		return fmt.Errorf("writeInitialPageData(): file.WriteAt: %w", err)
	}

	return nil
}
