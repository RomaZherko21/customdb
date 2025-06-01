package data

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
	"custom-database/internal/disk_manager/meta"
	"fmt"
)

// serializePage преобразует Page в []byte для записи на диск
func (fc *fileConnection) serializePage(page *Page) []byte {
	buffer := make([]byte, PAGE_SIZE)

	copy(buffer, fc.serializePageHeader(&page.Header))

	copy(buffer[PAGE_HEADER_SIZE:], fc.serializePageSlots(page.Slots))

	copy(buffer[PAGE_HEADER_SIZE+SLOTS_SPACE:], fc.serializePageData(page.Data))

	return buffer
}

// deserializePage восстанавливает Page из файла
func (fc *fileConnection) deserializePage(pageID uint32) *Page {
	pageHeader, err := fc.deserializePageHeader(pageID)
	if err != nil {
		return nil
	}

	slots, err := fc.deserializePageSlots(pageID, nil)
	if err != nil {
		return nil
	}

	pageData, err := fc.deserializePageData(pageID)
	if err != nil {
		return nil
	}

	return &Page{
		Header: *pageHeader,
		Slots:  slots,
		Data:   pageData,
	}
}

// serializePageHeader сериализует заголовок страницы
func (fc *fileConnection) serializePageHeader(pageHeader *PageHeader) []byte {
	buffer := make([]byte, PAGE_HEADER_SIZE)

	bs.WriteUint32(buffer, 0, pageHeader.PageId)
	bs.WriteUint16(buffer, PAGE_ID_SIZE, pageHeader.PageSize)

	return buffer
}

// deserializePageHeader десериализует заголовок страницы
func (fc *fileConnection) deserializePageHeader(pageID uint32) (*PageHeader, error) {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)

	pageData, err := fc.ReadFileRange(pageStartingPosition, pageStartingPosition+PAGE_HEADER_SIZE)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageHeader(): file.Read: %w", err)
	}

	pageHeader := &PageHeader{}
	pageHeader.PageId = bs.ReadUint32(pageData, 0)
	pageHeader.PageSize = bs.ReadUint16(pageData, PAGE_ID_SIZE)

	return pageHeader, nil
}

type interval struct {
	start uint32
	end   uint32
}

// serializePageSlots сериализует слоты страницы
func (fc *fileConnection) serializePageSlots(pageSlots []PageSlot) []byte {
	buffer := make([]byte, len(pageSlots)*ONE_SLOT_SIZE)

	for i, slot := range pageSlots {
		offset := i * ONE_SLOT_SIZE
		bs.WriteUint16(buffer, offset, slot.RowId)
		bs.WriteUint16(buffer, offset+SLOT_ROW_ID_SIZE, slot.Offset)
		bs.WriteUint16(buffer, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE, slot.Size)
		bs.WriteBool(buffer, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE, slot.IsDeleted)
	}

	return buffer
}

// deserializePageSlots десериализует слоты страницы
func (fc *fileConnection) deserializePageSlots(pageID uint32, interval *interval) ([]PageSlot, error) {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)

	start := pageStartingPosition + PAGE_HEADER_SIZE
	end := pageStartingPosition + PAGE_HEADER_SIZE + SLOTS_SPACE

	if interval != nil {
		start = interval.start
		end = interval.end
	}

	if (end-start)%ONE_SLOT_SIZE != 0 {
		return nil, fmt.Errorf("DeserializePageSlots(): end - start is not divisible by ONE_SLOT_SIZE")
	}
	slotsAmount := (end - start) / ONE_SLOT_SIZE

	pageData, err := fc.ReadFileRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageSlots(): file.Read: %w", err)
	}

	slots := make([]PageSlot, slotsAmount)
	for i := 0; i < int(slotsAmount); i++ {
		offset := i * ONE_SLOT_SIZE
		rowID := bs.ReadUint16(pageData, offset)
		// Если rowID != 0, значит слот содержит данные
		if rowID != 0 {
			slot := &PageSlot{
				RowId:     rowID,
				Offset:    bs.ReadUint16(pageData, offset+SLOT_ROW_ID_SIZE),
				Size:      bs.ReadUint16(pageData, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE),
				IsDeleted: bs.ReadBool(pageData, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE),
			}
			slots[i] = *slot
		}
	}

	return slots, nil
}

// serializePageData сериализует данные страницы
func (fc *fileConnection) serializePageData(pageData []byte) []byte {
	return pageData
}

// deserializePageData десериализует данные страницы
func (fc *fileConnection) deserializePageData(pageID uint32) ([]byte, error) {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)

	pageData, err := fc.ReadFileRange(pageStartingPosition+PAGE_HEADER_SIZE+SLOTS_SPACE, pageStartingPosition+PAGE_HEADER_SIZE+SLOTS_SPACE+DATA_SIZE)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageHeader(): file.Read: %w", err)
	}

	return pageData, nil
}

// deserializePageData десериализует данные страницы
func (fc *fileConnection) deserializeDataRow(pageID uint32, pageSlot *PageSlot, columns []*meta.Column) (*DataRow, error) {
	start := fc.CalculateDataRowPosition(pageID, pageSlot.Offset)
	end := start + uint32(pageSlot.Size)

	pageData, err := fc.ReadFileRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageHeader(): file.Read: %w", err)
	}

	offset := meta.NULL_BITMAP_SIZE

	row := make([]DataCell, len(columns))
	for i, column := range columns {
		columnValue, columnSize := meta.ConvertValueToType(pageData, offset, column.Type)
		row[i] = DataCell{
			Value:  columnValue,
			Type:   column.Type,
			IsNull: false,
		}
		offset += columnSize
	}

	return &DataRow{
		PageId: pageID,
		SlotId: pageSlot.RowId,
		Row:    row,
	}, nil
}
