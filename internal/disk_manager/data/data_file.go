package data

import (
	"fmt"
	"os"
	"path/filepath"

	bs "custom-database/internal/disk_manager/binary_serializer"
)

const (
	PAGE_SIZE = 4096 // 4KB

	MAX_SLOTS = 32 // Max slots on page

	PAGE_ID_SIZE     = 4
	PAGE_SIZE_SIZE   = 2
	PAGE_HEADER_SIZE = PAGE_ID_SIZE + PAGE_SIZE_SIZE + 2 // 2 bytes for padding

	SLOT_ROW_ID_SIZE     = 2
	SLOT_OFFSET_SIZE     = 2
	SLOT_SIZE_SIZE       = 2
	SLOT_IS_DELETED_SIZE = 1
	ONE_SLOT_SIZE        = SLOT_ROW_ID_SIZE + SLOT_OFFSET_SIZE + SLOT_SIZE_SIZE + SLOT_IS_DELETED_SIZE + 1

	SLOTS_SPACE = ONE_SLOT_SIZE * MAX_SLOTS                  // 8 * 32 = 256 bytes
	DATA_SIZE   = PAGE_SIZE - PAGE_HEADER_SIZE - SLOTS_SPACE // 4096 - 8 - 256 = 3808 bytes
)

func CreateDataFile(pageID uint32, filename string, filePath string) error {
	if _, err := os.Stat(filepath.Join(filePath, filename+".data")); err == nil {
		return fmt.Errorf("CreateDataFile(): table already exists: %w", err)
	}

	file, err := os.Create(filepath.Join(filePath, filename+".data"))
	if err != nil {
		return fmt.Errorf("CreateDataFile(): os.Create: %w", err)
	}
	defer file.Close()

	page := newPage(pageID)
	data := serializePage(page)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("CreateDataFile(): file.Write: %w", err)
	}

	return nil
}

func newPage(pageID uint32) *Page {
	return &Page{
		Header: PageHeader{
			PageId:   pageID,
			PageSize: PAGE_SIZE,
		},
		Slots: make([]PageSlot, MAX_SLOTS),
		Data:  make([]byte, DATA_SIZE),
	}
}

// serializePage преобразует Page в []byte для записи на диск
func serializePage(page *Page) []byte {
	// Выделяем буфер для всей страницы
	buffer := make([]byte, PAGE_SIZE)

	// 1. Сериализуем заголовок (первые PAGE_HEADER_SIZE байт)
	bs.WriteUint32(buffer, 0, page.Header.PageId)
	bs.WriteUint16(buffer, 4, page.Header.PageSize)

	// 2. Сериализуем слоты (следующие SLOTS_SPACE байт)
	slotsOffset := PAGE_HEADER_SIZE
	for i, slot := range page.Slots {
		offset := slotsOffset + (i * ONE_SLOT_SIZE)
		bs.WriteUint16(buffer, offset, slot.RowId)
		bs.WriteUint16(buffer, offset+4, slot.Offset)
		bs.WriteUint16(buffer, offset+8, slot.Size)
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
		Header: PageHeader{},
		Slots:  make([]PageSlot, MAX_SLOTS),
		Data:   make([]byte, DATA_SIZE),
	}

	// 1. Десериализуем заголовок
	page.Header.PageId = bs.ReadUint32(data, 0)
	page.Header.PageSize = bs.ReadUint16(data, 4)

	// 2. Десериализуем слоты
	slotsOffset := PAGE_HEADER_SIZE
	for i := 0; i < MAX_SLOTS; i++ {
		offset := slotsOffset + (i * ONE_SLOT_SIZE)
		rowID := bs.ReadUint16(data, offset)
		// Если rowID != 0, значит слот содержит данные
		if rowID != 0 {
			slot := &PageSlot{
				RowId:     rowID,
				Offset:    bs.ReadUint16(data, offset+4),
				Size:      bs.ReadUint16(data, offset+8),
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
