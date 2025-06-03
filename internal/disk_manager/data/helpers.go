package data

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
	"custom-database/internal/disk_manager/meta"
	"fmt"
	"io"
)

func (fc *fileConnection) CalculatePageStartingPosition(pageID uint32) uint32 {
	return (pageID - 1) * PAGE_SIZE
}

// CalculateDataRowPosition вычисляет позицию данных в строке
// Slot.offset - это смещение от начала страницы (страница = 4096 байт, offset = 4050, тогда позиция = 4050)
func (fc *fileConnection) CalculateDataRowPosition(pageID uint32, offset uint16) uint32 {
	pageStartingPosition := fc.CalculatePageStartingPosition(pageID)
	return pageStartingPosition + uint32(offset)
}

func (fc *fileConnection) ReadFileRange(start uint32, end uint32) ([]byte, error) {
	fc.file.Seek(int64(start), io.SeekStart)
	result := make([]byte, end-start)
	_, err := fc.file.Read(result)
	if err != nil {
		return nil, fmt.Errorf("ReadFileRange(): file.Read in range %d-%d: %w", start, end, err)
	}

	return result, nil
}

func CalculateDataRowSize(row []DataCell) uint32 {
	rowSize := meta.NULL_BITMAP_SIZE

	for _, cell := range row {
		if cell.IsNull {
			continue
		}

		if cell.Type == meta.TypeText {
			rowSize += len(cell.Value.(string)) + bs.TEXT_TYPE_HEADER
		} else {
			rowSize += meta.СalculateColumnSize(cell.Type)
		}
	}
	return uint32(rowSize)
}
