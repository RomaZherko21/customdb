package disk_manager

import (
	"fmt"
	"os"
)

const (
	META_FILE_SIZE = 4096 // 4KB
	MAX_COLUMNS    = 32   // Максимальное количество колонок в таблице
)

func createMetaFile(metaFile *MetaFile) {
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
	buffer := make([]byte, META_FILE_SIZE)

	// 1. Сериализуем имя таблицы
	offset := writeString(buffer, 0, metaFile.Name)

	// 2. Сериализуем количество колонок
	writeInt32(buffer, offset, int32(len(metaFile.Columns)))
	offset += 4

	// 3. Сериализуем null nullBitmap
	nullBitmap := uint32(0)
	for i := 0; i < len(metaFile.Columns); i++ {
		if metaFile.Columns[i].IsNullable {
			nullBitmap = setBit(nullBitmap, i)
		}
	}
	writeUint32(buffer, offset, nullBitmap)
	offset += 4

	// 4. Сериализуем каждую колонку
	for _, column := range metaFile.Columns {
		// [N байт] имя колонки
		offset += writeString(buffer, offset, column.Name)
		// [4 байт] тип данных (enum)
		writeInt32(buffer, offset, int32(column.Type))
		offset += 4
	}
	return buffer
}

func deserializeMetaFile(data []byte) *MetaFile {
	metaFile := &MetaFile{}

	// 1. Читаем имя таблицы
	fileName, offset := readString(data, 0)
	metaFile.Name = fileName

	// 2. Читаем количество колонок
	columnsCount := readInt32(data, offset)
	metaFile.Columns = make([]Column, columnsCount)
	offset += 4

	// 3. Читаем bitmap для nullable колонок
	nullBitmap := readUint32(data, offset)
	offset += 4

	// 4. Читаем информацию о колонках
	for i := 0; i < len(metaFile.Columns); i++ {
		columnName, columnNameOffset := readString(data, offset)
		columnType := ColumnType(readInt32(data, offset+columnNameOffset))
		offset += columnNameOffset + 4

		metaFile.Columns[i] = Column{
			Name:       columnName,
			Type:       columnType,
			IsNullable: getBit(nullBitmap, i),
		}
	}

	return metaFile
}

func setBit(bitmap uint32, position int) uint32 {
	return bitmap | (1 << position)
}

func getBit(bitmap uint32, position int) bool {
	return (bitmap & (1 << position)) != 0
}
