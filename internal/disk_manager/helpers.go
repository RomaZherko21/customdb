package disk_manager

// Вспомогательные функции для работы с байтами
// Little-endian порядок байтов

func writeUint8(buffer []byte, offset int, value uint8) {
	buffer[offset] = byte(value)
}

func readUint8(buffer []byte, offset int) uint8 {
	return uint8(buffer[offset])
}

func writeInt32(buffer []byte, offset int, value int32) {
	buffer[offset+3] = byte(value)
	buffer[offset+2] = byte(value >> 8)
	buffer[offset+1] = byte(value >> 16)
	buffer[offset] = byte(value >> 24)
}

func readInt32(buffer []byte, offset int) int32 {
	return int32(buffer[offset+3]) |
		int32(buffer[offset+2])<<8 |
		int32(buffer[offset+1])<<16 |
		int32(buffer[offset])<<24
}

func writeUint32(buffer []byte, offset int, value uint32) {
	buffer[offset+3] = byte(value)
	buffer[offset+2] = byte(value >> 8)
	buffer[offset+1] = byte(value >> 16)
	buffer[offset] = byte(value >> 24)
}

func readUint32(buffer []byte, offset int) uint32 {
	return uint32(buffer[offset+3]) |
		uint32(buffer[offset+2])<<8 |
		uint32(buffer[offset+1])<<16 |
		uint32(buffer[offset])<<24
}

func writeInt64(buffer []byte, offset int, value int64) {
	buffer[offset+7] = byte(value)
	buffer[offset+6] = byte(value >> 8)
	buffer[offset+5] = byte(value >> 16)
	buffer[offset+4] = byte(value >> 24)
	buffer[offset+3] = byte(value >> 32)
	buffer[offset+2] = byte(value >> 40)
	buffer[offset+1] = byte(value >> 48)
	buffer[offset] = byte(value >> 56)
}

func readInt64(buffer []byte, offset int) int64 {
	return int64(buffer[offset+7]) |
		int64(buffer[offset+6])<<8 |
		int64(buffer[offset+5])<<16 |
		int64(buffer[offset+4])<<24 |
		int64(buffer[offset+3])<<32 |
		int64(buffer[offset+2])<<40 |
		int64(buffer[offset+1])<<48 |
		int64(buffer[offset])<<56
}

func writeUint64(buffer []byte, offset int, value uint64) {
	buffer[offset+7] = byte(value)
	buffer[offset+6] = byte(value >> 8)
	buffer[offset+5] = byte(value >> 16)
	buffer[offset+4] = byte(value >> 24)
	buffer[offset+3] = byte(value >> 32)
	buffer[offset+2] = byte(value >> 40)
	buffer[offset+1] = byte(value >> 48)
	buffer[offset] = byte(value >> 56)
}

func readUint64(buffer []byte, offset int) uint64 {
	return uint64(buffer[offset+7]) |
		uint64(buffer[offset+6])<<8 |
		uint64(buffer[offset+5])<<16 |
		uint64(buffer[offset+4])<<24 |
		uint64(buffer[offset+3])<<32 |
		uint64(buffer[offset+2])<<40 |
		uint64(buffer[offset+1])<<48 |
		uint64(buffer[offset])<<56
}

// Записывает bool значение в 1 байт
func writeBool(buffer []byte, offset int, value bool) {
	if value {
		buffer[offset] = 1
	} else {
		buffer[offset] = 0
	}
}

// Читает bool значение из 1 байта
func readBool(buffer []byte, offset int) bool {
	return buffer[offset] != 0
}

const (
	TEXT_TYPE_HEADER = 4 // Размер длины строки в int32
)

// Записывает строку в буфер
// Формат: [длина строки (4 байта)][данные строки]
// Возвращает количество записанных байт
func writeString(buffer []byte, offset int, value string) int {
	// Записываем длину строки
	strLen := int32(len(value))
	writeInt32(buffer, offset, strLen)

	// Записываем сами данные
	copy(buffer[TEXT_TYPE_HEADER+offset:], []byte(value))

	// Возвращаем общее количество записанных байт
	return TEXT_TYPE_HEADER + int(strLen)
}

// Читает строку из буфера
// Возвращает прочитанную строку и количество прочитанных байт
func readString(buffer []byte, offset int) (string, int) {
	// Читаем длину строки
	strLen := readInt32(buffer, offset)

	// Читаем данные строки
	data := buffer[TEXT_TYPE_HEADER+offset : TEXT_TYPE_HEADER+offset+int(strLen)]

	// Возвращаем строку и общее количество прочитанных байт
	return string(data), TEXT_TYPE_HEADER + int(strLen)
}
