package data

import (
	"fmt"
	"io"
)

func (fc *fileConnection) CalculatePageStartingPosition(pageID uint32) uint32 {
	return (pageID - 1) * PAGE_SIZE
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
