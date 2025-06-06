package data

import (
	"fmt"
	"io"

	bs "custom-database/internal/disk_manager/binary_serializer"
	helpers "custom-database/internal/disk_manager/helpers"
)

func (fc *dataService) writeMetaData(tableName string, columns []Column) (*MetaData, int, error) {
	metaData := &MetaData{
		Name:      tableName,
		PageCount: 1,
		Columns:   columns,
	}

	data := serializeMetaData(metaData)

	_, err := fc.file.Write(data)
	if err != nil {
		return nil, 0, fmt.Errorf("WriteMetaData(): file.Write: %w", err)
	}

	fc.metaDataSpace = len(data)
	fc.meta = metaData

	return metaData, len(data), nil
}

func (fc *dataService) loadMetaData() (*MetaData, int, error) {
	data, err := io.ReadAll(fc.file)
	if err != nil {
		return nil, 0, fmt.Errorf("loadMetaData(): io.ReadAll: %w", err)
	}

	metaData, offset := deserializeMetaData(data)

	return metaData, offset, nil
}

// serializePage преобразует Page в []byte для записи на диск
func serializeMetaData(metaData *MetaData) []byte {
	// Выделяем буфер для всей страницы
	buffer := make([]byte, calculateFileSize(metaData))

	// 1. Сериализуем имя таблицы
	offset := bs.WriteString(buffer, 0, metaData.Name)

	// 2. Сериализуем количество страниц
	bs.WriteUint32(buffer, offset, metaData.PageCount)
	offset += PAGE_COUNT_SIZE

	// 3. Сериализуем количество колонок
	bs.WriteUint8(buffer, offset, uint8(len(metaData.Columns)))
	offset += COLUMN_COUNT_SIZE

	// 4. Сериализуем null nullBitmap
	nullBitmap := uint32(0)
	for i := 0; i < len(metaData.Columns); i++ {
		if metaData.Columns[i].IsNullable {
			nullBitmap = helpers.SetBit(nullBitmap, i)
		} else {
			nullBitmap = helpers.ClearBit(nullBitmap, i)
		}
	}
	bs.WriteUint32(buffer, offset, nullBitmap)
	offset += NULL_BITMAP_SIZE

	// 5. Сериализуем каждую колонку
	for _, column := range metaData.Columns {
		// [N байт] имя колонки
		offset += bs.WriteString(buffer, offset, column.Name)
		// [1 байт] тип данных (enum)
		bs.WriteUint8(buffer, offset, uint8(column.Type))
		offset += DATA_TYPE_SIZE
	}
	return buffer
}

func deserializeMetaData(data []byte) (*MetaData, int) {
	offset := 0
	metaData := &MetaData{}

	// 1. Читаем имя таблицы
	fileName, offset := bs.ReadString(data, offset)
	metaData.Name = fileName

	// 2. Читаем количество страниц
	metaData.PageCount = bs.ReadUint32(data, offset)
	offset += PAGE_COUNT_SIZE

	// 3. Читаем количество колонок
	columnsCount := bs.ReadUint8(data, offset)
	metaData.Columns = make([]Column, columnsCount)
	offset += COLUMN_COUNT_SIZE

	// 4. Читаем bitmap для nullable колонок
	nullBitmap := bs.ReadUint32(data, offset)
	offset += NULL_BITMAP_SIZE

	// 5. Читаем информацию о колонках
	for i := 0; i < len(metaData.Columns); i++ {
		columnName, columnNameOffset := bs.ReadString(data, offset)

		columnType := ColumnType(bs.ReadUint8(data, offset+columnNameOffset))
		offset += columnNameOffset + DATA_TYPE_SIZE

		metaData.Columns[i] = Column{
			Name:       columnName,
			Type:       columnType,
			IsNullable: helpers.GetBit(nullBitmap, i),
		}
	}

	return metaData, offset
}
