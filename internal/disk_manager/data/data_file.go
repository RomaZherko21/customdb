package data

import (
	"custom-database/internal/disk_manager/meta"
	"fmt"
	"os"
	"path/filepath"
)

type fileConnection struct {
	lastPageID uint32
	file       *os.File
	meta       *meta.MetaFile
}

const INITIAL_PAGE_ID = 1

func NewFileConnection(metaFile *meta.MetaFile, filePath string, isNewFile bool) (*fileConnection, error) {
	filePath = filepath.Join(filePath, metaFile.Name+".data")

	if isNewFile {
		fc := &fileConnection{
			lastPageID: INITIAL_PAGE_ID,
		}

		if err := fc.createDataFile(filePath); err != nil {
			return nil, fmt.Errorf("NewFileConnection(): CreateDataFile: %w", err)
		}
		return fc, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("NewFileConnection(): os.Open: %w", err)
	}

	return &fileConnection{
		lastPageID: INITIAL_PAGE_ID,
		file:       file,
		meta:       metaFile,
	}, nil
}

func (fc *fileConnection) Close() error {
	if fc.file == nil {
		return nil
	}

	return fc.file.Close()
}

type InsertRowResult struct {
	SlotID uint16
	PageID uint32
}

func (fc *fileConnection) InsertDataRow(row []DataCell) (*InsertRowResult, error) {
	pageHeaders, err := fc.deserializeAllPageHeaders(fc.meta.PageCount)
	if err != nil {
		return nil, fmt.Errorf("Insert(): deserializeAllPageHeaders: %w", err)
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
		// TODO: handle case when no free space is found
	}

	serializedRow := fc.serializeDataRow(row)

	slot, err := fc.InsertPageSlot(currentPageID, uint16(rowSize))
	if err != nil {
		return nil, fmt.Errorf("Insert(): InsertPageSlot: %w", err)
	}

	writeDataPosition := slot.Offset
	_, err = fc.file.WriteAt(serializedRow, int64(writeDataPosition))
	if err != nil {
		return nil, fmt.Errorf("Insert(): file.WriteAt: %w", err)
	}

	return &InsertRowResult{SlotID: slot.SlotId, PageID: currentPageID}, nil
}

func (fc *fileConnection) InsertPageSlot(pageID uint32, rowSize uint16) (*PageSlot, error) {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)

	slots, err := fc.deserializePageSlots(pageID, nil)
	if err != nil {
		return nil, fmt.Errorf("Insert(): deserializePageSlots: %w", err)
	}

	lastSlotOffset := PAGE_SIZE - rowSize
	if len(slots) > 0 {
		lastSlotOffset = slots[len(slots)-1].Offset - rowSize
	}

	newSlotID := len(slots) + 1

	slot := PageSlot{
		SlotId:    uint16(newSlotID),
		Offset:    uint16(lastSlotOffset - rowSize),
		RowSize:   rowSize,
		IsDeleted: false,
	}

	serializedSlot := fc.serializePageSlots([]PageSlot{slot})

	writeSlotPosition := pageStartingPosition + PAGE_HEADER_SIZE + uint32(len(slots)*ONE_SLOT_SIZE)
	fc.file.WriteAt(serializedSlot, int64(writeSlotPosition))

	return &slot, nil
}

func (fc *fileConnection) createDataFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("createDataFile(): table already exists: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("createDataFile(): os.Create: %w", err)
	}
	fc.file = file

	page := fc.newPage(INITIAL_PAGE_ID)
	data := fc.serializePage(page)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("createDataFile(): file.Write: %w", err)
	}

	return nil
}

func (fc *fileConnection) newPage(pageID uint32) *Page {
	return &Page{
		Header: PageHeader{
			PageId:    pageID,
			FreeSpace: DATA_SPACE,
		},
		Slots: make([]PageSlot, MAX_SLOTS),
		Data:  make([]DataRow, 0),
	}
}
