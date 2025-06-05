package data

import (
	bs "custom-database/internal/disk_manager/binary_serializer"
	"custom-database/internal/disk_manager/meta"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTableName = "test_table"
const testFilePath = "./test_data"

func TestSerializePageHeader(t *testing.T) {
	fileConnection := &fileConnection{}

	t.Run("Создание таблицы", func(t *testing.T) {
		// Создаем тестовую директорию
		err := os.MkdirAll(testFilePath, 0755)
		assert.NoError(t, err)

		// Создаем тестовые колонки
		columns := []meta.Column{
			{Name: "id", Type: meta.TypeInt32, IsNullable: false},
			{Name: "name", Type: meta.TypeText, IsNullable: true},
		}

		// Создаем метафайл
		metaFile, err := meta.CreateMetaFile(testTableName, columns, testFilePath)
		assert.NoError(t, err)

		// Создаем соединение с файлом
		fc, err := NewFileConnection(metaFile, testFilePath, true)
		assert.NoError(t, err)

		fileConnection = fc
	})

	t.Run("корректная сериализация заголовка страницы", func(t *testing.T) {
		header := &PageHeader{
			PageId:    1,
			FreeSpace: 100,
		}

		result := fileConnection.serializePageHeader(header)

		assert.Equal(t, PAGE_HEADER_SIZE, len(result))
		assert.Equal(t, uint32(1), bs.ReadUint32(result, 0))
		assert.Equal(t, uint16(100), bs.ReadUint16(result, PAGE_ID_SIZE))
	})

	t.Run("корректная десериализация заголовка страницы", func(t *testing.T) {
		pageID := uint32(1)
		result, err := fileConnection.deserializePageHeader(pageID)

		assert.NoError(t, err)
		assert.Equal(t, pageID, result.PageId)
		assert.Equal(t, uint16(DATA_SPACE), result.FreeSpace)
	})

	t.Run("корректная сериализация слотов", func(t *testing.T) {
		slots := []PageSlot{
			{SlotId: 1, Offset: 100, RowSize: 50, IsDeleted: false},
			{SlotId: 2, Offset: 200, RowSize: 75, IsDeleted: true},
		}

		result := fileConnection.serializePageSlots(slots)

		assert.Equal(t, len(slots)*ONE_SLOT_SIZE, len(result))
		assert.Equal(t, uint16(1), bs.ReadUint16(result, 0))
		assert.Equal(t, uint16(100), bs.ReadUint16(result, SLOT_ROW_ID_SIZE))
		assert.Equal(t, uint16(50), bs.ReadUint16(result, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE))
		assert.False(t, bs.ReadBool(result, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE))
	})

	t.Run("корректная десериализация слотов", func(t *testing.T) {
		pageID := uint32(1)
		result, err := fileConnection.deserializePageSlots(pageID, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, MAX_SLOTS, len(result))
	})

	t.Run("корректная сериализация данных страницы", func(t *testing.T) {
		pageData := []DataRow{
			{
				PageId: 1,
				SlotId: 1,
				Row: []DataCell{
					{Value: int32(111), Type: meta.TypeInt32, IsNull: false},
					{Value: "test111", Type: meta.TypeText, IsNull: false},
				},
			},
			{
				PageId: 1,
				SlotId: 2,
				Row: []DataCell{
					{Value: int32(222), Type: meta.TypeInt32, IsNull: false},
					{Value: "test222", Type: meta.TypeText, IsNull: false},
				},
			},
		}

		row1Size := int(CalculateDataRowSize(pageData[0].Row))

		result := fileConnection.serializePageData(pageData)
		assert.NotEmpty(t, result)

		assert.Equal(t, int32(111), bs.ReadInt32(result, meta.NULL_BITMAP_SIZE))
		assert.Equal(t, int32(222), bs.ReadInt32(result, row1Size+meta.NULL_BITMAP_SIZE))

		row1, _ := bs.ReadString(result, meta.NULL_BITMAP_SIZE+4)
		row2, _ := bs.ReadString(result, row1Size+meta.NULL_BITMAP_SIZE+4)
		assert.Equal(t, "test111", row1)
		assert.Equal(t, "test222", row2)
	})

	t.Run("корректная десериализация данных страницы", func(t *testing.T) {
		pageID := uint32(1)
		result, err := fileConnection.deserializePageData(pageID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result) // Новая страница должна быть пустой
	})

	t.Run("корректная сериализация строки данных", func(t *testing.T) {
		dataRow := []DataCell{
			{Value: int32(1), Type: meta.TypeInt32, IsNull: false},
			{Value: "test", Type: meta.TypeText, IsNull: true},
		}

		result := fileConnection.serializeDataRow(dataRow)

		assert.NotEmpty(t, result)
		assert.Equal(t, uint32(2), bs.ReadUint32(result, 0)) // nullBitmap для второй колонки
	})

	t.Run("корректная десериализация всех заголовков страниц", func(t *testing.T) {
		result, err := fileConnection.deserializeAllPageHeaders(1)

		assert.NoError(t, err)
		assert.Len(t, result, 1) // Начальная страница
		assert.Equal(t, uint32(1), result[0].PageId)
		assert.Equal(t, uint16(DATA_SPACE), result[0].FreeSpace)
	})

	t.Run("Закрытие файла", func(t *testing.T) {
		err := fileConnection.Close()
		assert.NoError(t, err)

		err = os.RemoveAll(testFilePath)
		assert.NoError(t, err)
	})
}
