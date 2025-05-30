package meta

import (
	"fmt"
	"os"

	bs "custom-database/internal/disk_manager/binary_serializer"
)

const (
	MAX_COLUMNS = 32 // Максимальное количество колонок в таблице

	NULL_BITMAP_SIZE  = 4 // Размер null_bitmap в uint32
	COLUMN_COUNT_SIZE = 1 // Размер количества колонок в uint8
	DATA_TYPE_SIZE    = 1 // Размер типа данных в uint8
)

func CreateMetaFile(metaFile *MetaFile) {
	file, err := os.Create(metaFile.Name + ".meta")
	if err != nil {
		fmt.Printf("Failed to create file: %v", err)
	}
	defer file.Close()

	data := serializeMetaFile(metaFile)

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("Failed to write page: %v", err)
	}
}

// serializePage преобразует Page в []byte для записи на диск
func serializeMetaFile(metaFile *MetaFile) []byte {
	// Выделяем буфер для всей страницы
	buffer := make([]byte, calculateFileSize(metaFile))

	// 1. Сериализуем имя таблицы
	offset := bs.WriteString(buffer, 0, metaFile.Name)

	// 2. Сериализуем количество колонок
	bs.WriteUint8(buffer, offset, uint8(len(metaFile.Columns)))
	offset += COLUMN_COUNT_SIZE

	// 3. Сериализуем null nullBitmap
	nullBitmap := uint32(0)
	for i := 0; i < len(metaFile.Columns); i++ {
		if metaFile.Columns[i].IsNullable {
			nullBitmap = setBit(nullBitmap, i)
		} else {
			nullBitmap = clearBit(nullBitmap, i)
		}
	}
	bs.WriteUint32(buffer, offset, nullBitmap)
	offset += NULL_BITMAP_SIZE

	// 4. Сериализуем каждую колонку
	for _, column := range metaFile.Columns {
		// [N байт] имя колонки
		offset += bs.WriteString(buffer, offset, column.Name)
		// [1 байт] тип данных (enum)
		bs.WriteUint8(buffer, offset, uint8(column.Type))
		offset += DATA_TYPE_SIZE
	}
	return buffer
}

func deserializeMetaFile(data []byte) *MetaFile {
	metaFile := &MetaFile{}

	// 1. Читаем имя таблицы
	fileName, offset := bs.ReadString(data, 0)
	metaFile.Name = fileName

	// 2. Читаем количество колонок
	columnsCount := bs.ReadUint8(data, offset)
	metaFile.Columns = make([]Column, columnsCount)
	offset += COLUMN_COUNT_SIZE

	// 3. Читаем bitmap для nullable колонок
	nullBitmap := bs.ReadUint32(data, offset)
	offset += NULL_BITMAP_SIZE

	// 4. Читаем информацию о колонках
	for i := 0; i < len(metaFile.Columns); i++ {
		columnName, columnNameOffset := bs.ReadString(data, offset)

		columnType := ColumnType(bs.ReadUint8(data, offset+columnNameOffset))
		offset += columnNameOffset + DATA_TYPE_SIZE

		metaFile.Columns[i] = Column{
			Name:       columnName,
			Type:       columnType,
			IsNullable: getBit(nullBitmap, i),
		}
	}

	return metaFile
}
