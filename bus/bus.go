package bus

import (
	"nes_emulator/loader"
	"nes_emulator/mlog"
)

const RamSize = 64 * 1024

type Bus struct {
	//Cpu *Device
	ram [RamSize]uint8
}

func New() *Bus {
	b := Bus{}
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

func (b *Bus) LoadNes(nes *loader.NesFile) {
	// Todo: load nes directly
}

func (b *Bus) LoadToRam(data []uint8, startAddr int) {
	for offset, dataByte := range data {
		addr := startAddr + offset
		if addr >= 0 && addr <= 0xFFFF {
			b.ram[addr] = dataByte
		} else {
			mlog.L.Fatal("LoadToRam memory out of range!")
		}
	}
}
