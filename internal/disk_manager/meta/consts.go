package meta

const (
	MAX_COLUMNS = 32 // Максимальное количество колонок в таблице

	NULL_BITMAP_SIZE  = 4 // Размер null_bitmap в uint32
	COLUMN_COUNT_SIZE = 1 // Размер количества колонок в uint8
	DATA_TYPE_SIZE    = 1 // Размер типа данных в uint8
)
