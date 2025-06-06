package data

import (
	"testing"

	bs "custom-database/internal/disk_manager/binary_serializer"
	helpers "custom-database/internal/disk_manager/helpers"

	"github.com/stretchr/testify/assert"
)

func TestSerializePageHeader(t *testing.T) {
	t.Run("Serialize PageHeader (PageId=1, FreeSpace=100, SlotsAmount=0)", func(t *testing.T) {
		header := &PageHeader{PageId: 1, FreeSpace: 100, SlotsAmount: 0}
		buf := serializePageHeader(header)
		assert.Equal(t, PAGE_HEADER_SIZE, len(buf))
		assert.Equal(t, uint32(1), bs.ReadUint32(buf, 0))
		assert.Equal(t, uint16(100), bs.ReadUint16(buf, PAGE_ID_SIZE))
		assert.Equal(t, uint16(0), bs.ReadUint16(buf, PAGE_ID_SIZE+PAGE_SLOTS_AMOUNT_SIZE))

		deserialized, err := deserializePageHeader(buf)
		assert.NoError(t, err)
		assert.Equal(t, header, deserialized)
	})

	t.Run("Serialize PageHeader (PageId=2, FreeSpace=200, SlotsAmount=5)", func(t *testing.T) {
		header := &PageHeader{PageId: 2, FreeSpace: 200, SlotsAmount: 5}
		buf := serializePageHeader(header)
		assert.Equal(t, PAGE_HEADER_SIZE, len(buf))
		assert.Equal(t, uint32(2), bs.ReadUint32(buf, 0))
		assert.Equal(t, uint16(200), bs.ReadUint16(buf, PAGE_ID_SIZE))
		assert.Equal(t, uint16(5), bs.ReadUint16(buf, PAGE_ID_SIZE+PAGE_SLOTS_AMOUNT_SIZE))

		deserialized, err := deserializePageHeader(buf)
		assert.NoError(t, err)
		assert.Equal(t, header, deserialized)
	})

	t.Run("Serialize PageHeader (PageId=0, FreeSpace=0, SlotsAmount=0)", func(t *testing.T) {
		header := &PageHeader{PageId: 0, FreeSpace: 0, SlotsAmount: 0}
		buf := serializePageHeader(header)
		assert.Equal(t, PAGE_HEADER_SIZE, len(buf))
		assert.Equal(t, uint32(0), bs.ReadUint32(buf, 0))
		assert.Equal(t, uint16(0), bs.ReadUint16(buf, PAGE_ID_SIZE))
		assert.Equal(t, uint16(0), bs.ReadUint16(buf, PAGE_ID_SIZE+PAGE_SLOTS_AMOUNT_SIZE))

		deserialized, err := deserializePageHeader(buf)
		assert.NoError(t, err)
		assert.Equal(t, header, deserialized)
	})
}

func TestSerializePageSlots(t *testing.T) {
	t.Run("Serialize empty slots array", func(t *testing.T) {
		slots := []PageSlot{}
		buf := serializePageSlots(slots)
		assert.Equal(t, 0, len(buf), "SerializePageSlots: ожидалась пустая длина буфера для пустого массива слотов")
	})

	t.Run("Serialize single slot", func(t *testing.T) {
		slots := []PageSlot{
			{
				SlotId:    1,
				Offset:    100,
				RowSize:   50,
				IsDeleted: false,
			},
		}
		buf := serializePageSlots(slots)
		assert.Equal(t, ONE_SLOT_SIZE, len(buf), "SerializePageSlots: ожидалась длина буфера %d, получено %d", ONE_SLOT_SIZE, len(buf))
		assert.Equal(t, uint16(1), bs.ReadUint16(buf, 0), "SerializePageSlots: ожидался SlotId=1, получено другое значение")
		assert.Equal(t, uint16(100), bs.ReadUint16(buf, SLOT_ROW_ID_SIZE), "SerializePageSlots: ожидался Offset=100, получено другое значение")
		assert.Equal(t, uint16(50), bs.ReadUint16(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE), "SerializePageSlots: ожидался RowSize=50, получено другое значение")
		assert.Equal(t, false, bs.ReadBool(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE), "SerializePageSlots: ожидался IsDeleted=false, получено другое значение")
	})

	t.Run("Serialize multiple slots", func(t *testing.T) {
		slots := []PageSlot{
			{
				SlotId:    1,
				Offset:    100,
				RowSize:   50,
				IsDeleted: false,
			},
			{
				SlotId:    2,
				Offset:    200,
				RowSize:   75,
				IsDeleted: true,
			},
		}
		buf := serializePageSlots(slots)
		assert.Equal(t, ONE_SLOT_SIZE*2, len(buf), "SerializePageSlots: ожидалась длина буфера %d, получено %d", ONE_SLOT_SIZE*2, len(buf))

		// Проверяем первый слот
		assert.Equal(t, uint16(1), bs.ReadUint16(buf, 0), "SerializePageSlots: ожидался SlotId=1, получено другое значение")
		assert.Equal(t, uint16(100), bs.ReadUint16(buf, SLOT_ROW_ID_SIZE), "SerializePageSlots: ожидался Offset=100, получено другое значение")
		assert.Equal(t, uint16(50), bs.ReadUint16(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE), "SerializePageSlots: ожидался RowSize=50, получено другое значение")
		assert.Equal(t, false, bs.ReadBool(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE), "SerializePageSlots: ожидался IsDeleted=false, получено другое значение")

		// Проверяем второй слот
		secondSlotOffset := ONE_SLOT_SIZE
		assert.Equal(t, uint16(2), bs.ReadUint16(buf, secondSlotOffset), "SerializePageSlots: ожидался SlotId=2, получено другое значение")
		assert.Equal(t, uint16(200), bs.ReadUint16(buf, secondSlotOffset+SLOT_ROW_ID_SIZE), "SerializePageSlots: ожидался Offset=200, получено другое значение")
		assert.Equal(t, uint16(75), bs.ReadUint16(buf, secondSlotOffset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE), "SerializePageSlots: ожидался RowSize=75, получено другое значение")
		assert.Equal(t, true, bs.ReadBool(buf, secondSlotOffset+SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE), "SerializePageSlots: ожидался IsDeleted=true, получено другое значение")
	})
}

func TestDeserializePageSlots(t *testing.T) {
	t.Run("Deserialize empty buffer", func(t *testing.T) {
		slots, err := deserializePageSlots([]byte{})
		assert.NoError(t, err, "DeserializePageSlots: ошибка при десериализации пустого буфера: %v", err)
		assert.Equal(t, 0, len(slots), "DeserializePageSlots: ожидался пустой массив слотов")
	})

	t.Run("Deserialize single slot", func(t *testing.T) {
		slots := []PageSlot{
			{
				SlotId:    1,
				Offset:    4050,
				RowSize:   50,
				IsDeleted: false,
			},
			{
				SlotId:    2,
				Offset:    4000,
				RowSize:   50,
				IsDeleted: false,
			},
		}

		buf := serializePageSlots(slots)

		result, err := deserializePageSlots(buf)
		assert.NoError(t, err, "DeserializePageSlots: ошибка при десериализации одного слота: %v", err)
		assert.Equal(t, slots, result, "DeserializePageSlots: ожидался массив из одного слота")
	})

	// Подтест: десериализация слота с нулевым SlotId (пустой слот)
	t.Run("Deserialize slot with zero SlotId", func(t *testing.T) {
		buf := make([]byte, ONE_SLOT_SIZE)
		bs.WriteUint16(buf, 0, 0)                                                  // SlotId = 0 (пустой слот)
		bs.WriteUint16(buf, SLOT_ROW_ID_SIZE, 100)                                 // Offset
		bs.WriteUint16(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE, 50)                 // RowSize
		bs.WriteBool(buf, SLOT_ROW_ID_SIZE+SLOT_OFFSET_SIZE+SLOT_SIZE_SIZE, false) // IsDeleted

		slots, err := deserializePageSlots(buf)
		assert.NoError(t, err, "DeserializePageSlots: ошибка при десериализации слота с нулевым SlotId: %v", err)
		assert.Equal(t, 0, len(slots), "DeserializePageSlots: ожидался пустой массив слотов для слота с нулевым SlotId")
	})
}

func TestSerializeDataRow(t *testing.T) {
	t.Run("Serialize row with all values", func(t *testing.T) {
		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
			{Value: "test", Type: TypeText, IsNull: false},
			{Value: uint32(25), Type: TypeUint32, IsNull: false},
			{Value: true, Type: TypeBoolean, IsNull: false},
		}

		buf := serializeDataRow(row)
		expectedSize := NULL_BITMAP_SIZE + // null bitmap
			4 + // int32
			(4 + len("test")) + // text (length + value)
			4 + // uint32
			1 // boolean
		assert.Equal(t, expectedSize, len(buf))

		// Проверяем null bitmap (все значения не null)
		nullBitmap := bs.ReadUint32(buf, 0)
		assert.Equal(t, uint32(0), nullBitmap)

		// Проверяем значения
		offset := NULL_BITMAP_SIZE
		assert.Equal(t, int32(1), bs.ReadInt32(buf, offset))
		offset += 4

		textValue, textLen := bs.ReadString(buf, offset)
		assert.Equal(t, "test", textValue)
		offset += textLen

		assert.Equal(t, uint32(25), bs.ReadUint32(buf, offset))
		offset += 4

		assert.Equal(t, true, bs.ReadBool(buf, offset))
	})

	t.Run("Serialize row with null values", func(t *testing.T) {
		row := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
			{Value: nil, Type: TypeText, IsNull: true},
			{Value: uint32(25), Type: TypeUint32, IsNull: false},
			{Value: nil, Type: TypeBoolean, IsNull: true},
		}

		buf := serializeDataRow(row)
		expectedSize := NULL_BITMAP_SIZE + // null bitmap
			4 + // int32
			4 + // uint32
			0 // null values are not serialized
		assert.Equal(t, expectedSize, len(buf))

		// Проверяем null bitmap (биты 1 и 3 установлены)
		nullBitmap := bs.ReadUint32(buf, 0)
		expectedBitmap := uint32(0)
		expectedBitmap = helpers.SetBit(expectedBitmap, 1) // второй бит (индекс 1)
		expectedBitmap = helpers.SetBit(expectedBitmap, 3) // четвертый бит (индекс 3)
		assert.Equal(t, expectedBitmap, nullBitmap)

		// Проверяем значения (только не null)
		offset := NULL_BITMAP_SIZE
		assert.Equal(t, int32(1), bs.ReadInt32(buf, offset))
		offset += 4

		assert.Equal(t, uint32(25), bs.ReadUint32(buf, offset))
	})

	t.Run("Serialize empty row", func(t *testing.T) {
		row := []DataCell{}
		buf := serializeDataRow(row)
		assert.Equal(t, NULL_BITMAP_SIZE, len(buf))
		assert.Equal(t, uint32(0), bs.ReadUint32(buf, 0))
	})
}

func TestDeserializeDataRow(t *testing.T) {
	columns := []Column{
		{Name: "id", Type: TypeInt32, IsNullable: false},
		{Name: "name", Type: TypeText, IsNullable: true},
		{Name: "age", Type: TypeUint32, IsNullable: false},
		{Name: "is_active", Type: TypeBoolean, IsNullable: true},
	}

	t.Run("Deserialize row with all values", func(t *testing.T) {
		// Создаем буфер с данными
		buf := make([]byte, NULL_BITMAP_SIZE+4+(4+len("test"))+4+1)

		// Записываем null bitmap (все значения не null)
		bs.WriteUint32(buf, 0, 0)

		// Записываем значения
		offset := NULL_BITMAP_SIZE
		bs.WriteInt32(buf, offset, 1)
		offset += 4

		bs.WriteString(buf, offset, "test")
		offset += 4 + len("test")

		bs.WriteUint32(buf, offset, 25)
		offset += 4

		bs.WriteBool(buf, offset, true)

		// Десериализуем
		row, err := deserializeDataRow(buf, columns)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(row))

		// Проверяем значения
		assert.Equal(t, int32(1), row[0].Value)
		assert.Equal(t, TypeInt32, row[0].Type)
		assert.False(t, row[0].IsNull)

		assert.Equal(t, "test", row[1].Value)
		assert.Equal(t, TypeText, row[1].Type)
		assert.False(t, row[1].IsNull)

		assert.Equal(t, uint32(25), row[2].Value)
		assert.Equal(t, TypeUint32, row[2].Type)
		assert.False(t, row[2].IsNull)

		assert.Equal(t, true, row[3].Value)
		assert.Equal(t, TypeBoolean, row[3].Type)
		assert.False(t, row[3].IsNull)
	})

	t.Run("Deserialize row with null values", func(t *testing.T) {
		// Создаем буфер с данными
		buf := make([]byte, NULL_BITMAP_SIZE+4+4) // только для не null значений

		// Записываем null bitmap (биты 1 и 3 установлены)
		nullBitmap := uint32(0)
		nullBitmap = helpers.SetBit(nullBitmap, 1)
		nullBitmap = helpers.SetBit(nullBitmap, 3)
		bs.WriteUint32(buf, 0, nullBitmap)

		// Записываем только не null значения
		offset := NULL_BITMAP_SIZE
		bs.WriteInt32(buf, offset, 1)
		offset += 4

		bs.WriteUint32(buf, offset, 25)

		// Десериализуем
		row, err := deserializeDataRow(buf, columns)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(row))

		// Проверяем значения
		assert.Equal(t, int32(1), row[0].Value)
		assert.Equal(t, TypeInt32, row[0].Type)
		assert.False(t, row[0].IsNull)

		assert.Nil(t, row[1].Value)
		assert.Equal(t, TypeText, row[1].Type)
		assert.True(t, row[1].IsNull)

		assert.Equal(t, uint32(25), row[2].Value)
		assert.Equal(t, TypeUint32, row[2].Type)
		assert.False(t, row[2].IsNull)

		assert.Nil(t, row[3].Value)
		assert.Equal(t, TypeBoolean, row[3].Type)
		assert.True(t, row[3].IsNull)
	})

	t.Run("Deserialize empty row", func(t *testing.T) {
		buf := make([]byte, NULL_BITMAP_SIZE)
		bs.WriteUint32(buf, 0, 0)

		row, err := deserializeDataRow(buf, []Column{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(row))
	})

	t.Run("Round trip serialization", func(t *testing.T) {
		columns := []Column{
			{Name: "id", Type: TypeInt32, IsNullable: false},
			{Name: "name", Type: TypeText, IsNullable: true},
			{Name: "age", Type: TypeUint32, IsNullable: true},
			{Name: "is_active", Type: TypeBoolean, IsNullable: true},
		}

		originalRow := []DataCell{
			{Value: int32(1), Type: TypeInt32, IsNull: false},
			{Value: "test", Type: TypeText, IsNull: false},
			{Value: nil, Type: TypeUint32, IsNull: true},
			{Value: true, Type: TypeBoolean, IsNull: false},
		}

		// Сериализуем
		buf := serializeDataRow(originalRow)

		// Десериализуем
		deserializedRow, err := deserializeDataRow(buf, columns)
		assert.NoError(t, err)
		assert.Equal(t, len(originalRow), len(deserializedRow))

		// Проверяем каждое значение
		for i := range originalRow {
			assert.Equal(t, originalRow[i].Value, deserializedRow[i].Value)
			assert.Equal(t, originalRow[i].Type, deserializedRow[i].Type)
			assert.Equal(t, originalRow[i].IsNull, deserializedRow[i].IsNull)
		}
	})
}
