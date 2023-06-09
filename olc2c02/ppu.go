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
	// OAM
	oamMem     [64]spriteObjEntry
	oamAddrBuf uint8
	// OAM drawing state
	nextLineScanlineSprites   [8]spriteObjEntry
	nextLineSpriteCount       uint8
	nextLineSpriteShiftPtrnLo [8]uint8
	nextLineSpriteShiftPtrnHi [8]uint8
	// Detection for sprite 0 collision
	spriteZeroHitPossible bool
	spriteZeroRendering   bool
	// Internal
	lookupPalette  [PreDefPaletteSize]color.RGBA
	scanLine       int16
	cycle          uint16
	frameCompleted bool
	hasNmi         bool
	// Screen display
	screenDisplayBuf    []uint8
	nametableDisplayBuf []uint8
	patternDisplayBuf   [2][]uint8
	width               int
	height              int
	// Registers
	regCtrl register
	regMask register
	regStat register
	// nametable scrolling, t = temp: https://www.nesdev.org/wiki/PPU_scrolling with loopy register
	regLoopyVram registerLoopy // this is for implementing loopy rendering! Not used in NES
	regLoopyTram registerLoopy
	scrollFineX  uint8
	// addr buffer (2 cycle to write 2 bytes to ppu, so buffer first one)
	addressLatch bool
	dataBuffer   uint8
	// Next render tile information (ref: ppu clock update diagram)
	bgNextTileId     uint8
	bgNextTileAttrib uint8
	bgNextTileLsb    uint8
	bgNextTileMsb    uint8
	// Shifter for scrolling
	bgShifterPatternLow  uint16
	bgShifterPatternHigh uint16
	bgShifterAttribLow   uint16
	bgShifterAttribHigh  uint16
}

func New() *Ppu {
	p := &Ppu{
		screenDisplayBuf:    make([]byte, NesDisplayWidth*NesDisplayHeight*4),
		nametableDisplayBuf: make([]byte, NesDisplayWidth*NesDisplayHeight*4),
		addressLatch:        false,
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
	mlog.L.Infof("Ppu is initialised with screen dim (%v, %v)", NesDisplayWidth, NesDisplayHeight)
	return p
}

func (p *Ppu) ConnectBus(b *bus.Bus) {
	p.b = b
}

func (p *Ppu) ConnectDisk(nesDisk *disk.NesDisk) {
	p.disk = nesDisk
}
