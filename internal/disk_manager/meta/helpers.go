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
		PAGE_COUNT_SIZE + // количество страниц в uint32
		COLUMN_COUNT_SIZE + // количество колонок в uint8
		NULL_BITMAP_SIZE + // null_bitmap size в uint32
		columnSize
}

func СalculateColumnSize(columnType ColumnType) int {
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
		return bs.ReadInt32(data, offset), СalculateColumnSize(columnType)
	case TypeInt64:
		return bs.ReadInt64(data, offset), СalculateColumnSize(columnType)
	case TypeUint32:
		return bs.ReadUint32(data, offset), СalculateColumnSize(columnType)
	case TypeUint64:
		return bs.ReadUint64(data, offset), СalculateColumnSize(columnType)
	case TypeBoolean:
		return bs.ReadBool(data, offset), СalculateColumnSize(columnType)
	case TypeText:
		return bs.ReadString(data, offset)
	}

	panic("ConvertValueToType(): unknown column type")
}

func ConvertValueToBuffer(columnType ColumnType, value interface{}) []byte {
	buffer1 := make([]byte, 1)
	buffer4 := make([]byte, 4)
	buffer8 := make([]byte, 8)

	switch columnType {
	case TypeInt32:
		bs.WriteInt32(buffer4, 0, value.(int32))
		return buffer4
	case TypeInt64:
		bs.WriteInt64(buffer8, 0, value.(int64))
		return buffer8
	case TypeUint32:
		bs.WriteUint32(buffer4, 0, value.(uint32))
		return buffer4
	case TypeUint64:
		bs.WriteUint64(buffer8, 0, value.(uint64))
		return buffer8
	case TypeBoolean:
		bs.WriteBool(buffer1, 0, value.(bool))
		return buffer1
	case TypeText:
		buffer := make([]byte, bs.TEXT_TYPE_HEADER+len(value.(string)))
		bs.WriteString(buffer, 0, value.(string))
		return buffer
	}

	panic("ConvertValueToBuffer(): unknown column type")
}
