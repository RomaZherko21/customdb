package meta

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
)

func calculateFileSize(metaFile *MetaFile) int {
	columnSize := 0
	for _, column := range metaFile.Columns {
		columnSize += bs.TEXT_TYPE_HEADER + len(column.Name) + DATA_TYPE_SIZE
	}

	return bs.TEXT_TYPE_HEADER + len(metaFile.Name) +
		COLUMN_COUNT_SIZE + // количество колонок в uint8
		NULL_BITMAP_SIZE + // null_bitmap size в uint32
		columnSize
}

func setBit(bitmap uint32, position int) uint32 {
	return bitmap | (1 << position)
}

func clearBit(bitmap uint32, position int) uint32 {
	return bitmap &^ (1 << position)
}

func getBit(bitmap uint32, position int) bool {
	return (bitmap & (1 << position)) != 0
}
