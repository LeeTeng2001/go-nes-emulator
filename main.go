package main

import (
	"nes_emulator/bus"
	"nes_emulator/cpu6502"
)

func main() {
	// TODO: Main entry
	newBus := bus.New()
	cpu := cpu6502.New()
	cpu.ConnectBus(newBus)
	for {
		cpu.Clock()
	}
}
