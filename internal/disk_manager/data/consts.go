package data

// Meta data
const (
	PAGE_COUNT_SIZE   = 4 // Размер количества страниц в uint32
	NULL_BITMAP_SIZE  = 4 // Размер null_bitmap в uint32
	COLUMN_COUNT_SIZE = 1 // Размер количества колонок в uint8
	DATA_TYPE_SIZE    = 1 // Размер типа данных в uint8
)

// Page data
const (
	PAGE_SIZE = 4096 // 4KB

	MAX_SLOTS = 32 // Max slots on page

	PAGE_ID_SIZE           = 4
	PAGE_SIZE_SIZE         = 2
	PAGE_SLOTS_AMOUNT_SIZE = 2
	PAGE_HEADER_SIZE       = PAGE_ID_SIZE + PAGE_SIZE_SIZE + PAGE_SLOTS_AMOUNT_SIZE

	SLOT_ROW_ID_SIZE     = 2
	SLOT_OFFSET_SIZE     = 2
	SLOT_SIZE_SIZE       = 2
	SLOT_IS_DELETED_SIZE = 1
	ONE_SLOT_SIZE        = SLOT_ROW_ID_SIZE + SLOT_OFFSET_SIZE + SLOT_SIZE_SIZE + SLOT_IS_DELETED_SIZE + 1

	SLOTS_SPACE = ONE_SLOT_SIZE * MAX_SLOTS                  // 8 * 32 = 256 bytes
	DATA_SPACE  = PAGE_SIZE - PAGE_HEADER_SIZE - SLOTS_SPACE // 4096 - 8 - 256 = 3834 bytes
)

const INITIAL_PAGE_ID = 1
