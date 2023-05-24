package disk

import "testing"

func TestLoadNes(t *testing.T) {
	nesFile := New("tests/nestest.nes")
	nesFile.PrintInfo()
}
