package helpers

func SetBit(bitmap uint32, position int) uint32 {
	return bitmap | (1 << position)
}

func ClearBit(bitmap uint32, position int) uint32 {
	return bitmap &^ (1 << position)
}

func GetBit(bitmap uint32, position int) bool {
	return (bitmap & (1 << position)) != 0
}
