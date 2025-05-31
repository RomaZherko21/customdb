package data

import (
	"fmt"
	"os"

	bs "custom-database/internal/disk_manager/binary_serializer"
)

const (
	PAGE_SIZE        = 4096 // 4KB
	ONE_SLOT_SIZE    = 16   // 16 bytes (row_id + offset + size + is_deleted + padding)
	PAGE_HEADER_SIZE = 128  // 128 bytes (page_id + page_size + for other possible fields)
	MAX_SLOTS        = 32   // Max slots on page

	SLOTS_SPACE = ONE_SLOT_SIZE * MAX_SLOTS                  // 16 * 32 = 512 bytes
	DATA_SIZE   = PAGE_SIZE - PAGE_HEADER_SIZE - SLOTS_SPACE // 4096 - 128 - 512 = 3456 bytes
)

func newPage(pageID int32) *Page {
	return &Page{
		Header: &PageHeader{
			PageId:   pageID,
			PageSize: int32(PAGE_SIZE),
		},
		Slots: make([]PageSlot, MAX_SLOTS),
		Data:  make([]byte, DATA_SIZE),
	}
}

func CreateDataFile(filename string) {
	file, err := os.Create(filename + ".data")
	if err != nil {
		fmt.Printf("Failed to create file: %v", err)
	}
	defer file.Close()

	page := newPage(4)
	data := serializePage(page)

	_, err = file.Write(data)
	if err != nil {
		fmt.Printf("Failed to write page: %v", err)
	}
}

// serializePage преобразует Page в []byte для записи на диск
func serializePage(page *Page) []byte {
	// Выделяем буфер для всей страницы
	buffer := make([]byte, PAGE_SIZE)

	// 1. Сериализуем заголовок (первые PAGE_HEADER_SIZE байт)
	bs.WriteInt32(buffer, 0, page.Header.PageId)
	bs.WriteInt32(buffer, 4, page.Header.PageSize)

	// 2. Сериализуем слоты (следующие SLOTS_SPACE байт)
	slotsOffset := PAGE_HEADER_SIZE
	for i, slot := range page.Slots {
		offset := slotsOffset + (i * ONE_SLOT_SIZE)
		bs.WriteInt32(buffer, offset, slot.RowId)
		bs.WriteInt32(buffer, offset+4, slot.Offset)
		bs.WriteInt32(buffer, offset+8, slot.Size)
		bs.WriteBool(buffer, offset+12, slot.IsDeleted)
	}

	// 3. Копируем данные (оставшиеся DATA_SIZE байт)
	dataOffset := PAGE_HEADER_SIZE + SLOTS_SPACE
	copy(buffer[dataOffset:], page.Data)

	return buffer
}

// DeserializePage восстанавливает Page из []byte прочитанных с диска
func DeserializePage(data []byte) *Page {
	page := &Page{
		Header: &PageHeader{},
		Slots:  make([]PageSlot, MAX_SLOTS),
		Data:   make([]byte, DATA_SIZE),
	}

	// 1. Десериализуем заголовок
	page.Header.PageId = bs.ReadInt32(data, 0)
	page.Header.PageSize = bs.ReadInt32(data, 4)

	// 2. Десериализуем слоты
	slotsOffset := PAGE_HEADER_SIZE
	for i := 0; i < MAX_SLOTS; i++ {
		offset := slotsOffset + (i * ONE_SLOT_SIZE)
		rowID := bs.ReadInt32(data, offset)
		// Если rowID != 0, значит слот содержит данные
		if rowID != 0 {
			slot := &PageSlot{
				RowId:     rowID,
				Offset:    bs.ReadInt32(data, offset+4),
				Size:      bs.ReadInt32(data, offset+8),
				IsDeleted: bs.ReadBool(data, offset+12),
			}
			page.Slots[i] = *slot
		}
	}

	// 3. Копируем данные
	dataOffset := PAGE_HEADER_SIZE + SLOTS_SPACE
	copy(page.Data, data[dataOffset:dataOffset+DATA_SIZE])

	return page
}
