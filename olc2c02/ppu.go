package olc2c02

import (
	"image/color"
	"math"
	"nes_emulator/bus"
	"nes_emulator/disk"
)

// PPU memory map: https://www.nesdev.org/wiki/PPU_memory_map
// Ppu has 16kb address space completely separate from cpu bus.
// It has ability to address 4 nametable but it stored 2 nametable
// Has pallet information

// Compile time interface check
var _ bus.PpuDevice = (*Ppu)(nil)

const (
	NameTableSize     = 0x0400
	PreDefPaletteSize = 0x40
)

type Ppu struct {
	b    *bus.Bus
	disk *disk.NesDisk
	// VRAM
	tableName    [2][NameTableSize]uint8
	tablePalette [32]uint8
	// Internal
	lookupPalette [PreDefPaletteSize]color.Color
	scanLine      uint16
	cycle         uint16
}

func New() *Ppu {
	p := &Ppu{}
	p.initLookupPalette()
	return p
}

func (p *Ppu) ConnectBus(b *bus.Bus) {
	p.b = b
}

func (p *Ppu) ConnectDisk(nesDisk *disk.NesDisk) {
	p.disk = nesDisk
}

func (p *Ppu) Clock() {
	// RN is just random noise
	// Fake some noise for now
	//sprScreen->SetPixel(cycle - 1, scanline, palScreen[(rand() % 2) ? 0x3F : 0x30]);

	// Advance renderer - it never stops
	// TODO: Explain the numbers
	p.cycle++
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanLine++
		if p.scanLine >= 261 {
			p.scanLine = math.MaxUint16
		}
	}
}
