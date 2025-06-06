package data

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
	"custom-database/internal/disk_manager/helpers"
	"fmt"
)

// serializePage преобразует Page в []byte для записи на диск
func serializePage(page *Page) []byte {
	buffer := make([]byte, PAGE_SIZE)

	copy(buffer, serializePageHeader(&page.Header))

	copy(buffer[PAGE_HEADER_SIZE:], serializePageSlots(page.Slots))

	return buffer
}

// serializePageHeader сериализует заголовок страницы
func serializePageHeader(pageHeader *PageHeader) []byte {
	buffer := make([]byte, PAGE_HEADER_SIZE)

	bs.WriteUint32(buffer, 0, pageHeader.PageId)
	bs.WriteUint16(buffer, PAGE_ID_SIZE, pageHeader.FreeSpace)
	bs.WriteUint16(buffer, PAGE_ID_SIZE+PAGE_SLOTS_AMOUNT_SIZE, pageHeader.SlotsAmount)

	return buffer
}

// deserializePageHeader десериализует заголовок страницы
func deserializePageHeader(pageData []byte) (*PageHeader, error) {
	pageHeader := &PageHeader{}
	pageHeader.PageId = bs.ReadUint32(pageData, 0)
	pageHeader.FreeSpace = bs.ReadUint16(pageData, PAGE_ID_SIZE)
	pageHeader.SlotsAmount = bs.ReadUint16(pageData, PAGE_ID_SIZE+PAGE_SLOTS_AMOUNT_SIZE)

	return pageHeader, nil
}

// serializePageSlots сериализует слоты страницы
func serializePageSlots(pageSlots []PageSlot) []byte {
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
func deserializePageSlots(data []byte) ([]PageSlot, error) {
	if len(data)%ONE_SLOT_SIZE != 0 {
		return nil, fmt.Errorf("DeserializePageSlots(): data length is not divisible by ONE_SLOT_SIZE")
	}

	slotsAmount := len(data) / ONE_SLOT_SIZE
	slots := make([]PageSlot, 0)

	for i := 0; i < slotsAmount; i++ {
		offset := i * ONE_SLOT_SIZE
		rowID := bs.ReadUint16(data, offset)
		// Если rowID != 0, значит слот содержит данные
		if rowID != 0 {
			slot := &PageSlot{
				SlotId:    rowID,
				Offset:    bs.ReadUint16(data, offset+SLOT_ROW_ID_SIZE),
				RowSize:   bs.ReadUint16(data, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE),
				IsDeleted: bs.ReadBool(data, offset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE),
			}
			slots = append(slots, *slot)
		}
	}

	return slots, nil
}

// serializeDataRow сериализует данные строки
func serializeDataRow(dataRow []DataCell) []byte {
	rowSize := CalculateDataRowSize(dataRow)
	buffer := make([]byte, rowSize)

	nullBitmap := uint32(0)
	for i, cell := range dataRow {
		if cell.IsNull {
			nullBitmap = helpers.SetBit(nullBitmap, i)
		}
	}

	bs.WriteUint32(buffer, 0, nullBitmap)
	offset := NULL_BITMAP_SIZE

	for _, cell := range dataRow {
		if cell.IsNull {
			continue
		}

		byteValue := ConvertValueToBuffer(cell.Type, cell.Value)
		copy(buffer[offset:], byteValue)
		offset += len(byteValue)
	}

	return buffer
}

// deserializePageData десериализует данные страницы
func deserializeDataRow(dataRow []byte, columns []Column) ([]DataCell, error) {
	nullBitmap := bs.ReadUint32(dataRow, 0)
	offset := NULL_BITMAP_SIZE

	row := make([]DataCell, len(columns))
	for i, column := range columns {
		// Если колонка nullable и в nullBitmap на этой позиции стоит 1, то значение null
		if column.IsNullable && helpers.GetBit(nullBitmap, i) {
			row[i] = DataCell{
				Value:  nil,
				Type:   column.Type,
				IsNull: true,
			}
			continue
		}

		columnValue, columnSize := ConvertValueToType(dataRow, offset, column.Type)
		row[i] = DataCell{
			Value:  columnValue,
			Type:   column.Type,
			IsNull: false,
		}
		offset += columnSize
	}

	return row, nil
}
