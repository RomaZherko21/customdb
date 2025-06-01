package helpers

import (
	"fmt"
	"strings"
	"testing"
)

func TestBitOperations(t *testing.T) {
	t.Run("Установка, очистка и проверка бита", func(t *testing.T) {
		var bitmap uint32 = 0

		// Тест установки бита
		positions := []int{0, 1, 3, 31} // Проверяем разные позиции, включая граничные
		for _, pos := range positions {
			bitmap = SetBit(bitmap, pos)
			if !GetBit(bitmap, pos) {
				t.Errorf("Бит в позиции %d должен быть установлен", pos)
			}
			fmt.Printf("После установки бита %d: %s (число: %d)\n", pos, bitsToString(bitmap), bitmap)
		}

		// Тест очистки бита
		for _, pos := range positions {
			bitmap = ClearBit(bitmap, pos)
			if GetBit(bitmap, pos) {
				t.Errorf("Бит в позиции %d должен быть очищен", pos)
			}
			fmt.Printf("После очистки бита %d: %s (число: %d)\n", pos, bitsToString(bitmap), bitmap)
		}
	})

	t.Run("Множественные операции с битами", func(t *testing.T) {
		var bitmap uint32 = 0

		// Устанавливаем несколько битов
		bitmap = SetBit(bitmap, 1)
		bitmap = SetBit(bitmap, 3)
		bitmap = SetBit(bitmap, 5)

		// Проверяем установленные биты
		if !GetBit(bitmap, 1) || !GetBit(bitmap, 3) || !GetBit(bitmap, 5) {
			t.Error("Не все биты были корректно установлены")
		}

		// Проверяем неустановленные биты
		if GetBit(bitmap, 0) || GetBit(bitmap, 2) || GetBit(bitmap, 4) {
			t.Error("Обнаружены неожиданно установленные биты")
		}

		// Очищаем один бит и проверяем остальные
		bitmap = ClearBit(bitmap, 3)
		if !GetBit(bitmap, 1) || GetBit(bitmap, 3) || !GetBit(bitmap, 5) {
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
