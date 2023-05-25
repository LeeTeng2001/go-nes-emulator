package olc2c02

import (
	"image/color"
	"math/rand"
	"nes_emulator/bus"
	"nes_emulator/disk"
	"nes_emulator/mlog"
)

// PPU memory map, must read: https://www.nesdev.org/wiki/PPU_memory_map
// Ppu has 16kb address space completely separate from cpu bus.
// It has ability to address 4 nametable but it stored 2 nametable
// Has pallet information, all of the 4th pallette reflects the first background color

// Compile time interface check
var _ bus.PpuDevice = (*Ppu)(nil)

const (
	NesDisplayWidth   = 256
	NesDisplayHeight  = 240
	NameTableSize     = 0x0400
	PatternTableSize  = 128 * 128
	PreDefPaletteSize = 0x40
	PaletteRamOffset  = 0x3F00
)

type Ppu struct {
	b    *bus.Bus
	disk *disk.NesDisk
	// VRAM
	tableName     [2][NameTableSize]uint8
	tablePatterns [2][PatternTableSize]uint8
	tablePalette  [32]uint8
	// Internal
	lookupPalette  [PreDefPaletteSize]color.RGBA
	scanLine       int16
	cycle          uint16
	frameCompleted bool
	hasNmi         bool
	// Screen display
	screenDisplayBuf  []uint8
	patternDisplayBuf [2][]uint8
	width             int
	height            int
	// Registers
	regCtrl register
	regMask register
	regStat register
	// addr buffer (2 cycle to write 2 bytes to ppu, so buffer first one)
	hasWriteBuffer bool
	addrCombined   uint16
	dataBuffer     uint8
}

func New() *Ppu {
	p := &Ppu{
		screenDisplayBuf: make([]byte, NesDisplayWidth*NesDisplayHeight*4),
		hasWriteBuffer:   false,
	}
	p.patternDisplayBuf[0] = make([]byte, PatternTableSize*4)
	p.patternDisplayBuf[1] = make([]byte, PatternTableSize*4)

	// Initialise screenDisplayBuf with random black white pixel
	for i := 0; i < NesDisplayWidth*NesDisplayHeight*4; i += 4 {
		if rand.Float32() < 0.5 {
			p.screenDisplayBuf[i] = 255
			p.screenDisplayBuf[i+1] = 255
			p.screenDisplayBuf[i+2] = 255
		} else {
			p.screenDisplayBuf[i] = 0
			p.screenDisplayBuf[i+1] = 0
			p.screenDisplayBuf[i+2] = 0
		}
		p.screenDisplayBuf[i+3] = 255
	}

	p.initLookupPalette()
	mlog.L.Infof("Ppu is initialised with screenDisplayBuf dim (%v, %v)", NesDisplayWidth, NesDisplayHeight)
	return p
}

func (p *Ppu) ConnectBus(b *bus.Bus) {
	p.b = b
}

func (p *Ppu) ConnectDisk(nesDisk *disk.NesDisk) {
	p.disk = nesDisk
}

func (p *Ppu) CheckNmiAndTurnOff() bool {
	if p.hasNmi {
		p.hasNmi = false
		return true
	}
	return false
}

func (p *Ppu) FrameCompleteAndTurnOff() bool {
	if p.frameCompleted {
		p.frameCompleted = false
		return true
	}
	return false
}

func (p *Ppu) Clock() {
	// Check if we leave vertical blanking
	if p.scanLine == -1 && p.cycle == 1 {
		p.regStat.SetFlag(regStatVertZeroBlank, false)
	}

	// Check if we reach the first out of range scanline and emit NMI if enabled
	if p.scanLine == 241 && p.cycle == 1 {
		p.regStat.SetFlag(regStatVertZeroBlank, true)
		if p.regCtrl.GetFlag(regCtrlEnableNMI) {
			p.hasNmi = true
		}
	}

	// Fake some noise for now
	if p.cycle < 256 && p.scanLine >= 0 && p.scanLine < 240 {
		i := 4 * (int(p.cycle) + int(p.scanLine)*256)
		if rand.Float32() < 0.5 {
			p.screenDisplayBuf[i] = 255
			p.screenDisplayBuf[i+1] = 255
			p.screenDisplayBuf[i+2] = 255
		} else {
			p.screenDisplayBuf[i] = 0
			p.screenDisplayBuf[i+1] = 0
			p.screenDisplayBuf[i+2] = 0
		}
	}

	// Advance. why 341 and 261?
	// The actual dim is 256x240 but the scanline is bigger than that 341x261
	// The unseen scanline (down) is known as the vertical blanking period!
	p.cycle++
	p.frameCompleted = false
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanLine++
		if p.scanLine >= 261 {
			p.scanLine = -1 // reset to one line before the actual area
			p.frameCompleted = true
		}
	}
}
