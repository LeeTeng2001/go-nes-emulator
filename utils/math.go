package utils

// FlipByte reverse the order of the byte, ex: 0b11100000 becomes 0b00000111
func FlipByte(b uint8) uint8 {
	// https://stackoverflow.com/a/2602885
	b = (b&0xF0)>>4 | (b&0x0F)<<4
	b = (b&0xCC)>>2 | (b&0x33)<<2
	b = (b&0xAA)>>1 | (b&0x55)<<1
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
