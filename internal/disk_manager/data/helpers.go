package data

import (
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
		return nil, fmt.Errorf("ReadFileRange(): file.Read: %w", err)
	}

	return result, nil
}
