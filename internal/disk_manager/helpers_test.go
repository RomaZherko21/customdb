package disk_manager

import (
	"math"
	"testing"
)

// Вычисляет размер, необходимый для записи строки
func getStringSize(value string) int {
	return 4 + len(value) // 4 байта на длину + сама строка
}

func TestInt32Operations(t *testing.T) {
	t.Run("Positive Numbers", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(300000)

		writeInt32(buffer, 0, value)
		result := readInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Negative Numbers", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(-300000)

		writeInt32(buffer, 0, value)
		result := readInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Min Int32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(math.MinInt32)

		writeInt32(buffer, 0, value)
		result := readInt32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Max Int32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := int32(math.MaxInt32)

		writeInt32(buffer, 0, value)
		result := readInt32(buffer, 0)

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

		writeUint32(buffer, 0, value)
		result := readUint32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Max Uint32", func(t *testing.T) {
		buffer := make([]byte, 4)
		value := uint32(math.MaxUint32)

		writeUint32(buffer, 0, value)
		result := readUint32(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestInt64Operations(t *testing.T) {
	t.Run("Positive Numbers", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := int64(math.MaxInt64) // MaxInt64

		writeInt64(buffer, 0, value)
		result := readInt64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})

	t.Run("Negative Numbers", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := int64(math.MinInt64) // MinInt64

		writeInt64(buffer, 0, value)
		result := readInt64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestUint64Operations(t *testing.T) {
	t.Run("Large Number", func(t *testing.T) {
		buffer := make([]byte, 8)
		value := uint64(math.MaxUint64) // MaxUint64

		writeUint64(buffer, 0, value)
		result := readUint64(buffer, 0)

		t.Logf("Число %d в байтах: % x", value, buffer)
		if result != value {
			t.Errorf("got %d, want %d", result, value)
		}
	})
}

func TestBooleanOperations(t *testing.T) {
	t.Run("Write and Read True", func(t *testing.T) {
		buffer := make([]byte, 1)
		writeBool(buffer, 0, true)

		t.Logf("True в байтах: %08b", buffer[0])

		result := readBool(buffer, 0)
		if !result {
			t.Error("got false, want true")
		}
	})

	t.Run("Write and Read False", func(t *testing.T) {
		buffer := make([]byte, 1)
		writeBool(buffer, 0, false)

		t.Logf("False в байтах: %08b", buffer[0])

		result := readBool(buffer, 0)
		if result {
			t.Error("got true, want false")
		}
	})
}

func TestStringOperations(t *testing.T) {
	t.Run("Empty String", func(t *testing.T) {
		value := ""
		buffer := make([]byte, getStringSize(value))

		written := writeString(buffer, 0, value)
		got, read := readString(buffer, 0)

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
		buffer := make([]byte, getStringSize(value))

		written := writeString(buffer, 0, value)
		got, read := readString(buffer, 0)

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
		buffer := make([]byte, getStringSize(value))

		written := writeString(buffer, 0, value)
		got, read := readString(buffer, 0)

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
		buffer := make([]byte, getStringSize(value))

		written := writeString(buffer, 0, value)
		got, read := readString(buffer, 0)

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
		writeInt32(buffer, 0, 12345)
		writeString(buffer, 4, "Hello")
		writeBool(buffer, 13, true)
		writeUint64(buffer, 14, 98765)

		// Читаем и проверяем
		if readInt32(buffer, 0) != 12345 {
			t.Error("Failed to read Int32")
		}

		if str, _ := readString(buffer, 4); str != "Hello" {
			t.Error("Failed to read String")
		}

		if !readBool(buffer, 13) {
			t.Error("Failed to read Bool")
		}

		if readUint64(buffer, 14) != 98765 {
			t.Error("Failed to read Uint64")
		}

		t.Logf("Буфер с разными типами данных: % x", buffer[:30])
	})
}

func TestBoundaryValues(t *testing.T) {
	t.Run("Zero Values", func(t *testing.T) {
		buffer := make([]byte, 100)

		writeInt32(buffer, 0, 0)
		writeUint32(buffer, 4, 0)
		writeInt64(buffer, 8, 0)
		writeUint64(buffer, 16, 0)
		writeString(buffer, 24, "")
		writeBool(buffer, 28, false)

		if readInt32(buffer, 0) != 0 {
			t.Error("Failed to handle zero Int32")
		}
		if readUint32(buffer, 4) != 0 {
			t.Error("Failed to handle zero Uint32")
		}
		if readInt64(buffer, 8) != 0 {
			t.Error("Failed to handle zero Int64")
		}
		if readUint64(buffer, 16) != 0 {
			t.Error("Failed to handle zero Uint64")
		}
		if str, _ := readString(buffer, 24); str != "" {
			t.Error("Failed to handle empty string")
		}
		if readBool(buffer, 28) {
			t.Error("Failed to handle false boolean")
		}
	})
}
