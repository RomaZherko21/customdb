package meta

import (
	"fmt"
	"strings"
	"testing"

	bs "custom-database/internal/disk_manager/binary_serializer"
)

func TestCalculateFileSize(t *testing.T) {
	tests := []struct {
		name     string
		metaFile *MetaFile
		want     int
	}{
		{
			name: "Пустой файл без колонок",
			metaFile: &MetaFile{
				Name:    "test",
				Columns: []Column{},
			},
			want: bs.TEXT_TYPE_HEADER + 4 + // длина имени "test"
				COLUMN_COUNT_SIZE + // размер для количества колонок
				NULL_BITMAP_SIZE, // размер для null bitmap
		},
		{
			name: "Файл с одной колонкой",
			metaFile: &MetaFile{
				Name: "test",
				Columns: []Column{
					{Name: "col1"},
				},
			},
			want: bs.TEXT_TYPE_HEADER + 4 + // длина имени "test"
				COLUMN_COUNT_SIZE + // размер для количества колонок
				NULL_BITMAP_SIZE + // размер для null bitmap
				bs.TEXT_TYPE_HEADER + 4 + // длина имени колонки "col1"
				DATA_TYPE_SIZE, // размер для типа данных
		},
		{
			name: "Файл с несколькими колонками",
			metaFile: &MetaFile{
				Name: "test",
				Columns: []Column{
					{Name: "col1"},
					{Name: "column2"},
					{Name: "col3"},
				},
			},
			want: bs.TEXT_TYPE_HEADER + 4 + // длина имени "test"
				COLUMN_COUNT_SIZE + // размер для количества колонок
				NULL_BITMAP_SIZE + // размер для null bitmap
				(bs.TEXT_TYPE_HEADER + 4 + DATA_TYPE_SIZE) + // первая колонка
				(bs.TEXT_TYPE_HEADER + 7 + DATA_TYPE_SIZE) + // вторая колонка
				(bs.TEXT_TYPE_HEADER + 4 + DATA_TYPE_SIZE), // третья колонка
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateFileSize(tt.metaFile)
			if got != tt.want {
				t.Errorf("calculateFileSize() = %v, хотим %v", got, tt.want)
			}
		})
	}
}

func TestBitOperations(t *testing.T) {
	t.Run("Установка, очистка и проверка бита", func(t *testing.T) {
		var bitmap uint32 = 0

		// Тест установки бита
		positions := []int{0, 1, 3, 31} // Проверяем разные позиции, включая граничные
		for _, pos := range positions {
			bitmap = setBit(bitmap, pos)
			if !getBit(bitmap, pos) {
				t.Errorf("Бит в позиции %d должен быть установлен", pos)
			}
			fmt.Printf("После установки бита %d: %s (число: %d)\n", pos, bitsToString(bitmap), bitmap)
		}

		// Тест очистки бита
		for _, pos := range positions {
			bitmap = clearBit(bitmap, pos)
			if getBit(bitmap, pos) {
				t.Errorf("Бит в позиции %d должен быть очищен", pos)
			}
			fmt.Printf("После очистки бита %d: %s (число: %d)\n", pos, bitsToString(bitmap), bitmap)
		}
	})

	t.Run("Множественные операции с битами", func(t *testing.T) {
		var bitmap uint32 = 0

		// Устанавливаем несколько битов
		bitmap = setBit(bitmap, 1)
		bitmap = setBit(bitmap, 3)
		bitmap = setBit(bitmap, 5)

		// Проверяем установленные биты
		if !getBit(bitmap, 1) || !getBit(bitmap, 3) || !getBit(bitmap, 5) {
			t.Error("Не все биты были корректно установлены")
		}

		// Проверяем неустановленные биты
		if getBit(bitmap, 0) || getBit(bitmap, 2) || getBit(bitmap, 4) {
			t.Error("Обнаружены неожиданно установленные биты")
		}

		// Очищаем один бит и проверяем остальные
		bitmap = clearBit(bitmap, 3)
		if !getBit(bitmap, 1) || getBit(bitmap, 3) || !getBit(bitmap, 5) {
			t.Error("Очистка одного бита повлияла на другие биты")
		}
	})
}

// Функция для отображения битов uint32
func bitsToString(n uint32) string {
	var bits []string
	for i := 31; i >= 0; i-- {
		if i > 0 && i%8 == 0 {
			bits = append(bits, " ")
		}
		if (n & (1 << i)) != 0 {
			bits = append(bits, "1")
		} else {
			bits = append(bits, "0")
		}
	}
	return strings.Join(bits, "")
}
