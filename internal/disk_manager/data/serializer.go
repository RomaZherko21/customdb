package data

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
	"custom-database/internal/disk_manager/helpers"
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

	dataRows, err := fc.deserializePageData(pageID)
	if err != nil {
		return nil
	}

	return &Page{
		Header: *pageHeader,
		Slots:  slots,
		Data:   dataRows,
	}
}

// serializePageHeader сериализует заголовок страницы
func (fc *fileConnection) serializePageHeader(pageHeader *PageHeader) []byte {
	buffer := make([]byte, PAGE_HEADER_SIZE)

	bs.WriteUint32(buffer, 0, pageHeader.PageId)
	bs.WriteUint16(buffer, PAGE_ID_SIZE, pageHeader.FreeSpace)

	return buffer
}

// deserializePageHeader десериализует заголовок страницы
func (fc *fileConnection) deserializePageHeader(pageID uint32) (*PageHeader, error) {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)

	pageData, err := fc.ReadFileRange(pageStartingPosition, pageStartingPosition+PAGE_HEADER_SIZE)
	if err != nil {
		return nil, fmt.Errorf("deserializePageHeader(): fc.ReadFileRange: %w", err)
	}

	pageHeader := &PageHeader{}
	pageHeader.PageId = bs.ReadUint32(pageData, 0)
	pageHeader.FreeSpace = bs.ReadUint16(pageData, PAGE_ID_SIZE)

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
		bs.WriteUint16(buffer, offset, slot.SlotId)
		bs.WriteUint16(buffer, offset+SLOT_ROW_ID_SIZE, slot.Offset)
		bs.WriteUint16(buffer, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE, slot.RowSize)
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
				SlotId:    rowID,
				Offset:    bs.ReadUint16(pageData, offset+SLOT_ROW_ID_SIZE),
				RowSize:   bs.ReadUint16(pageData, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE),
				IsDeleted: bs.ReadBool(pageData, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE),
			}
			slots[i] = *slot
		}
	}

	return slots, nil
}

// serializePageData сериализует данные страницы
func (fc *fileConnection) serializePageData(pageData []DataRow) []byte {
	buffer := make([]byte, 0)

	for _, row := range pageData {
		buffer = append(buffer, fc.serializeDataRow(row.Row)...)
	}

	return buffer
}

// deserializePageData десериализует данные страницы
func (fc *fileConnection) deserializePageData(pageID uint32) ([]DataRow, error) {
	slots, err := fc.deserializePageSlots(pageID, nil)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageData(): deserializePageSlots: %w", err)
	}

	result := make([]DataRow, 0)
	for _, slot := range slots {
		if slot.SlotId != 0 {
			row, err := fc.deserializeDataRow(pageID, &slot)
			if err != nil {
				return nil, fmt.Errorf("DeserializePageData(): deserializeDataRow: %w", err)
			}
			result = append(result, *row)
		}
	}

	return result, nil
}

// serializeDataRow сериализует данные строки
func (fc *fileConnection) serializeDataRow(dataRow []DataCell) []byte {
	rowSize := CalculateDataRowSize(dataRow)
	buffer := make([]byte, rowSize)

	nullBitmap := uint32(0)
	for i, cell := range dataRow {
		if cell.IsNull {
			nullBitmap = helpers.SetBit(nullBitmap, i)
		}
	}

	bs.WriteUint32(buffer, 0, nullBitmap)
	offset := meta.NULL_BITMAP_SIZE

	for _, cell := range dataRow {
		if cell.IsNull {
			continue
		}

		byteValue := meta.ConvertValueToBuffer(cell.Type, cell.Value)
		copy(buffer[offset:], byteValue)
		offset += len(byteValue)
	}

	return buffer
}

// deserializePageData десериализует данные страницы
func (fc *fileConnection) deserializeDataRow(pageID uint32, pageSlot *PageSlot) (*DataRow, error) {
	start := fc.CalculateDataRowPosition(pageID, pageSlot.Offset)
	end := start + uint32(pageSlot.RowSize)

	pageData, err := fc.ReadFileRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("DeserializePageHeader(): file.Read: %w", err)
	}

	nullBitmap := bs.ReadUint32(pageData, 0)
	offset := meta.NULL_BITMAP_SIZE

	row := make([]DataCell, len(fc.meta.Columns))
	for i, column := range fc.meta.Columns {
		// Если колонка nullable и в nullBitmap на этой позиции стоит 1, то значение null
		if column.IsNullable && helpers.GetBit(nullBitmap, i) {
			row[i] = DataCell{
				Value:  nil,
				Type:   column.Type,
				IsNull: true,
			}
			continue
		}

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
		SlotId: pageSlot.SlotId,
		Row:    row,
	}, nil
}

func (fc *fileConnection) deserializeAllPageHeaders(pageCount uint32) ([]*PageHeader, error) {
	result := make([]*PageHeader, 0)

	for i := INITIAL_PAGE_ID; i <= int(pageCount); i++ {
		pageHeader, err := fc.deserializePageHeader(uint32(i))
		if err != nil {
			return nil, fmt.Errorf("DeserializePagesHeader(): deserializePageHeader: %w", err)
		}
		result = append(result, pageHeader)
	}

	return result, nil
}
