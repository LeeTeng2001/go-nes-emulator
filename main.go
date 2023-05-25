package main

import (
	"nes_emulator/nes"
)

func main() {
	//g := nes.New("olc2c02/tests/color_test.nes")
	g := nes.New("cpu6502/tests/nestest.nes")
	g.Run()
}
