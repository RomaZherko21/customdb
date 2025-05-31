package binary_serializer

import (
	"math"
	"testing"
)

// Вычисляет размер, необходимый для записи строки
func GetStringSize(value string) int {
	return 4 + len(value) // 4 байта на длину + сама строка
}

func TestInt32Operations(t *testing.T) {
	t.Run("Positive Numbers", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(300000)

		WriteInt32(buffer, 0, value)
		result := ReadInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Negative Numbers", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(-300000)

		WriteInt32(buffer, 0, value)
		result := ReadInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Min Int32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(math.MinInt32)

		WriteInt32(buffer, 0, value)
		result := ReadInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Max Int32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(math.MaxInt32)

		WriteInt32(buffer, 0, value)
		result := ReadInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestUint32Operations(t *testing.T) {
	t.Run("Small Number", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := uint32(300000)

		WriteUint32(buffer, 0, value)
		result := ReadUint32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Max Uint32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := uint32(math.MaxUint32)

		WriteUint32(buffer, 0, value)
		result := ReadUint32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestInt64Operations(t *testing.T) {
	t.Run("Positive Numbers", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := int64(math.MaxInt64)

		WriteInt64(buffer, 0, value)
		result := ReadInt64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Negative Numbers", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := int64(math.MinInt64)

		WriteInt64(buffer, 0, value)
		result := ReadInt64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestUint64Operations(t *testing.T) {
	t.Run("Large Number", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := uint64(math.MaxUint64)

		WriteUint64(buffer, 0, value)
		result := ReadUint64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestBooleanOperations(t *testing.T) {
	t.Run("Write and Read True", func(t *testing.T) {
		buffer := make([]byte, 1)
		WriteBool(buffer, 0, true)

		t.Logf("True в байтах: %08b", buffer[0])

		result := ReadBool(buffer, 0)
		if !result {
			t.Error("got false, want true")
		}
	})

	t.Run("Write and Read False", func(t *testing.T) {
		buffer := make([]byte, 1)
		WriteBool(buffer, 0, false)

		t.Logf("False в байтах: %08b", buffer[0])

		result := ReadBool(buffer, 0)
		if result {
			t.Error("got true, want false")
		}
	})
}

func TestStringOperations(t *testing.T) {
	t.Run("Empty String", func(t *testing.T) {
		value := ""
		buffer := make([]byte, GetStringSize(value))

		written := WriteString(buffer, 0, value)
		got, read := ReadString(buffer, 0)

		t.Logf("Пустая строка в байтах: % x", buffer)

		if written != read {
			t.Errorf("written %d bytes but read %d", written, read)
		}
		if got != value {
			t.Errorf("got %q, want %q", got, value)
		}
	})

	t.Run("ASCII String", func(t *testing.T) {
		value := "Hello, World!"
		buffer := make([]byte, GetStringSize(value))

		written := WriteString(buffer, 0, value)
		got, read := ReadString(buffer, 0)

		t.Logf("ASCII строка в байтах: % x", buffer)

		if written != read {
			t.Errorf("written %d bytes but read %d", written, read)
		}
		if got != value {
			t.Errorf("got %q, want %q", got, value)
		}
	})

	t.Run("UTF-8 String", func(t *testing.T) {
		value := "Привет, 世界!"
		buffer := make([]byte, GetStringSize(value))

		written := WriteString(buffer, 0, value)
		got, read := ReadString(buffer, 0)

		t.Logf("UTF-8 строка в байтах: % x", buffer)

		if written != read {
			t.Errorf("written %d bytes but read %d", written, read)
		}
		if got != value {
			t.Errorf("got %q, want %q", got, value)
		}
	})

	t.Run("String with Special Characters", func(t *testing.T) {
		value := "Tab\t Newline\n Quote\" Slash\\"
		buffer := make([]byte, GetStringSize(value))

		written := WriteString(buffer, 0, value)
		got, read := ReadString(buffer, 0)

		t.Logf("Строка со спецсимволами в байтах: % x", buffer)

		if written != read {
			t.Errorf("written %d bytes but read %d", written, read)
		}
		if got != value {
			t.Errorf("got %q, want %q", got, value)
		}
	})
}

func TestOffsetOperations(t *testing.T) {
	t.Run("Multiple Values at Different Offsets", func(t *testing.T) {
		buffer := make([]byte, 100)

		// Записываем разные типы данных с разными смещениями
		WriteInt32(buffer, 0, 12345)
		WriteString(buffer, 4, "Hello")
		WriteBool(buffer, 13, true)
		WriteUint64(buffer, 14, 98765)

		// Читаем и проверяем
		if ReadInt32(buffer, 0) != 12345 {
			t.Error("Failed to read Int32")
		}

		if str, _ := ReadString(buffer, 4); str != "Hello" {
			t.Error("Failed to read String")
		}

		if !ReadBool(buffer, 13) {
			t.Error("Failed to read Bool")
		}

		if ReadUint64(buffer, 14) != 98765 {
			t.Error("Failed to read Uint64")
		}

		t.Logf("Буфер с разными типами данных: % x", buffer[:30])
	})
}

func TestBoundaryValues(t *testing.T) {
	t.Run("Zero Values", func(t *testing.T) {
		buffer := make([]byte, 100)

		WriteInt32(buffer, 0, 0)
		WriteUint32(buffer, 4, 0)
		WriteInt64(buffer, 8, 0)
		WriteUint64(buffer, 16, 0)
		WriteString(buffer, 24, "")
		WriteBool(buffer, 28, false)

		if ReadInt32(buffer, 0) != 0 {
			t.Error("Failed to handle zero Int32")
		}
		if ReadUint32(buffer, 4) != 0 {
			t.Error("Failed to handle zero Uint32")
		}
		if ReadInt64(buffer, 8) != 0 {
			t.Error("Failed to handle zero Int64")
		}
		if ReadUint64(buffer, 16) != 0 {
			t.Error("Failed to handle zero Uint64")
		}
		if str, _ := ReadString(buffer, 24); str != "" {
			t.Error("Failed to handle empty string")
		}
		if ReadBool(buffer, 28) {
			t.Error("Failed to handle false boolean")
		}
	})
}
