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

func calculateColumnSize(columnType ColumnType) int {
	switch columnType {
	case TypeInt32:
		return 4
	case TypeInt64:
		return 8
	case TypeUint32:
		return 4
	case TypeUint64:
		return 8
	case TypeBoolean:
		return 1
	case TypeText:
		return 0
	}

	panic("CalculateColumnSize(): unknown column type")
}

func ConvertValueToType(data []byte, offset int, columnType ColumnType) (interface{}, int) {
	switch columnType {
	case TypeInt32:
		return bs.ReadInt32(data, offset), calculateColumnSize(columnType)
	case TypeInt64:
		return bs.ReadInt64(data, offset), calculateColumnSize(columnType)
	case TypeUint32:
		return bs.ReadUint32(data, offset), calculateColumnSize(columnType)
	case TypeUint64:
		return bs.ReadUint64(data, offset), calculateColumnSize(columnType)
	case TypeBoolean:
		return bs.ReadBool(data, offset), calculateColumnSize(columnType)
	case TypeText:
		return bs.ReadString(data, offset)
	}

	panic("ConvertValueToType(): unknown column type")
}
