package disk_manager

type ColumnType int32

const (
	TypeInt32 ColumnType = iota
	TypeInt64
	TypeUint32
	TypeUint64
	TypeBoolean
	TypeText
)

type Column struct {
	Name       string
	Type       ColumnType
	IsNullable bool
}

type MetaFile struct {
	Name    string
	Columns []Column
}

// [4 байта] длина имени таблицы
// [N байт] имя таблицы
// [4 байта] количество колонок
// [2 байта] null bitmap (1 бит на колонку)
// Для каждой колонки:
//   [4 байта] длина имени колонки
//   [N байт] имя колонки
//   [1 байт] тип данных (enum)
