package disk_manager

import (
	"testing"
)

func TestReadInt32LargeNumber(t *testing.T) {
	// Число 300,000 в Little-endian формате
	buffer := []byte{
		0x00, // buffer[0] (старший байт):  00000000
		0x04, // buffer[1]:                 00000100
		0x93, // buffer[2]:                 10010011
		0xE0, // buffer[3] (младший байт):  11100000
	}

	// Показываем каждый байт
	t.Log("Представление числа 300,000 в памяти (Little-endian):")
	t.Logf("Байт 3 (младший): %08b (0xE0) << 0  = %032b", buffer[3], int32(buffer[3]))
	t.Logf("Байт 2:           %08b (0x93) << 8  = %032b", buffer[2], int32(buffer[2])<<8)
	t.Logf("Байт 1:           %08b (0x04) << 16 = %032b", buffer[1], int32(buffer[1])<<16)
	t.Logf("Байт 0 (старший): %08b (0x00) << 24 = %032b", buffer[0], int32(buffer[0])<<24)

	// Читаем значение
	value := readInt32(buffer, 0)

	// Показываем результат операции побитового ИЛИ (|)
	t.Log("\nСборка числа из байтов:")
	t.Logf("Байт 3: %032b", int32(buffer[3]))
	t.Logf("Байт 2: %032b", int32(buffer[2])<<8)
	t.Logf("Байт 1: %032b", int32(buffer[1])<<16)
	t.Logf("Байт 0: %032b", int32(buffer[0])<<24)
	t.Logf("Результат (OR): %032b = %d", value, value)

	if value != 300000 {
		t.Errorf("readInt32() = %d, want 300000", value)
	}
}

func TestWriteInt32LargeNumber(t *testing.T) {
	// Создаем буфер для записи
	buffer := make([]byte, 4)

	// Записываем число 300,000
	writeInt32(buffer, 0, 300000)

	// Ожидаемые значения байтов
	expected := []byte{
		0x00, // 00000000
		0x04, // 00000100
		0x93, // 10010011
		0xE0, // 11100000
	}

	// Проверяем каждый байт
	t.Log("Проверка записи числа 300,000:")
	for i, b := range buffer {
		t.Logf("Байт %d: получили %08b (0x%02x), ожидали %08b (0x%02x)",
			i, b, b, expected[i], expected[i])
		if b != expected[i] {
			t.Errorf("byte[%d] = %02x, want %02x", i, b, expected[i])
		}
	}
}

func TestBooleanOperations(t *testing.T) {
	tests := []struct {
		name  string
		value bool
	}{
		{"True value", true},
		{"False value", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем буфер для одного байта
			buffer := make([]byte, 1)

			// Записываем значение
			writeBool(buffer, 0, tt.value)

			// Проверяем записанный байт
			expectedByte := byte(0)
			if tt.value {
				expectedByte = 1
			}
			if buffer[0] != expectedByte {
				t.Errorf("writeBool() wrote %d, want %d", buffer[0], expectedByte)
			}

			// Читаем значение обратно
			got := readBool(buffer, 0)
			if got != tt.value {
				t.Errorf("readBool() = %v, want %v", got, tt.value)
			}

			// Выводим для отладки
			t.Logf("Значение %v записано как: %08b", tt.value, buffer[0])
		})
	}
}

// Тест на чтение нестандартных значений (любое ненулевое значение считается true)
func TestReadBoolNonStandardValues(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{"Zero is false", 0x00, false},
		{"One is true", 0x01, true},
		{"Any non-zero is true (0xFF)", 0xFF, true},
		{"Any non-zero is true (0x80)", 0x80, true},
		{"Any non-zero is true (0x7F)", 0x7F, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := []byte{tt.input}
			got := readBool(buffer, 0)
			if got != tt.expected {
				t.Errorf("readBool() = %v, want %v for input byte 0x%02x", got, tt.expected, tt.input)
			}
		})
	}
}

func TestStringOperations(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Пустая строка", ""},
		{"Короткая строка", "Hello"},
		{"Русский текст", "Привет, мир!"},
		{"Длинная строка", "This is a longer string with special chars: !@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем буфер достаточного размера
			size := getStringSize(tt.value)
			buffer := make([]byte, size)

			// Записываем строку
			written := writeString(buffer, 0, tt.value)
			if written != size {
				t.Errorf("writeString() wrote %d bytes, want %d", written, size)
			}

			// Читаем строку обратно
			got, read := readString(buffer, 0)
			if read != size {
				t.Errorf("readString() read %d bytes, want %d", read, size)
			}
			if got != tt.value {
				t.Errorf("readString() = %q, want %q", got, tt.value)
			}

			// Выводим для отладки
			t.Logf("Строка: %q", tt.value)
			t.Logf("Размер: %d байт (4 байта длина + %d байт данные)", size, len(tt.value))
			t.Logf("Буфер: % x", buffer)
		})
	}
}

func TestCharacterEncoding(t *testing.T) {
	// Тестируем английскую букву 'A' и русскую букву 'П'
	value := "AП"

	// Создаем буфер
	size := getStringSize(value)
	buffer := make([]byte, size)

	// Записываем строку
	writeString(buffer, 0, value)

	// Анализируем каждый байт
	t.Log("Анализ кодирования строки 'AП':")

	// Длина строки (первые 4 байта)
	t.Logf("Длина строки (4 байта): %08b %08b %08b %08b",
		buffer[0], buffer[1], buffer[2], buffer[3])

	// Английская буква 'A'
	t.Log("\nАнглийская буква 'A':")
	t.Logf("Один байт:    %08b", buffer[4])
	t.Logf("Десятичное:   %d", buffer[4])
	t.Logf("Hex:          0x%02X", buffer[4])

	// Русская буква 'П'
	t.Log("\nРусская буква 'П' (2 байта UTF-8):")
	t.Logf("Первый байт:  %08b", buffer[5])
	t.Logf("Второй байт:  %08b", buffer[6])
	t.Logf("Десятичные:   %d %d", buffer[5], buffer[6])
	t.Logf("Hex:          0x%02X 0x%02X", buffer[5], buffer[6])

	// Объяснение структуры UTF-8 для русской буквы
	t.Log("\nСтруктура UTF-8 для 'П':")
	t.Log("1) Первый байт начинается с '110' - указывает на 2-байтовую последовательность")
	t.Log("2) Второй байт начинается с '10' - указывает на продолжение последовательности")

	// Проверяем, что можем прочитать обратно
	got, _ := readString(buffer, 0)
	if got != value {
		t.Errorf("readString() = %q, want %q", got, value)
	}
}
