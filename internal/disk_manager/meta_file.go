package disk_manager

import (
	"fmt"
	"os"
)

const (
	META_FILE_SIZE = 4096 // 4KB
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

	// 3. Сериализуем каждую колонку
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

	fileName, offset := readString(data, 0)
	metaFile.Name = fileName

	columnsCount := readInt32(data, offset)
	metaFile.Columns = make([]Column, columnsCount)
	offset += 4

	for i := 0; i < len(metaFile.Columns); i++ {
		columnName, columnNameOffset := readString(data, offset)
		columnType := ColumnType(readInt32(data, offset+columnNameOffset))
		offset += columnNameOffset + 4
		metaFile.Columns[i] = Column{Name: columnName, Type: columnType}
	}

	return metaFile
}
