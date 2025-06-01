package meta

import (
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
