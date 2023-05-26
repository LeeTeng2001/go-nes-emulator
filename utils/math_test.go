package utils

import "testing"

func TestFlipByte(t *testing.T) {
	byte1 := uint8(0b1111_0000)
	byte2 := uint8(0b0011_1000)
	byte3 := uint8(0b0100_1001)
	if FlipByte(byte1) != 0b0000_1111 {
		t.Fatalf("Incorrect result")
	}
	if FlipByte(byte2) != 0b0001_1100 {
		t.Fatalf("Incorrect result")
	}
	if FlipByte(byte3) != 0b1001_0010 {
		t.Fatalf("Incorrect result")
	}
}
