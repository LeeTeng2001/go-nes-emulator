package bus

import "nes_emulator/cpu6502"

const RamSize = 64 * 1024

type Bus struct {
	cpu *cpu6502.Cpu
	ram [RamSize]uint8
}

func New() *Bus {
	b := Bus{
		cpu: cpu6502.New(),
	}
	b.cpu.ConnectBus(&b)
	return &b
}

func (b *Bus) Write(addr uint16, data uint8) {
	if addr >= 0 && addr <= 0xFFFF {
		b.ram[addr] = data
	}
}

func (b *Bus) Read(addr uint16) (data uint8) {
	if addr >= 0 && addr <= 0xFFFF {
		return b.ram[addr]
	}
	return 0
}
