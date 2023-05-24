package olc2c02

import "image/color"

// For more palette: https://www.nesdev.org/wiki/PPU_palettes

func (p *Ppu) initLookupPalette() {
	p.lookupPalette[0x00] = color.RGBA{R: 84, G: 84, B: 84, A: 255}
	p.lookupPalette[0x01] = color.RGBA{R: 0, G: 30, B: 116, A: 255}
	p.lookupPalette[0x02] = color.RGBA{R: 8, G: 16, B: 144, A: 255}
	p.lookupPalette[0x03] = color.RGBA{R: 48, G: 0, B: 136, A: 255}
	p.lookupPalette[0x04] = color.RGBA{R: 68, G: 0, B: 100, A: 255}
	p.lookupPalette[0x05] = color.RGBA{R: 92, G: 0, B: 48, A: 255}
	p.lookupPalette[0x06] = color.RGBA{R: 84, G: 4, B: 0, A: 255}
	p.lookupPalette[0x07] = color.RGBA{R: 60, G: 24, B: 0, A: 255}
	p.lookupPalette[0x08] = color.RGBA{R: 32, G: 42, B: 0, A: 255}
	p.lookupPalette[0x09] = color.RGBA{R: 8, G: 58, B: 0, A: 255}
	p.lookupPalette[0x0A] = color.RGBA{R: 0, G: 64, B: 0, A: 255}
	p.lookupPalette[0x0B] = color.RGBA{R: 0, G: 60, B: 0, A: 255}
	p.lookupPalette[0x0C] = color.RGBA{R: 0, G: 50, B: 60, A: 255}
	p.lookupPalette[0x0D] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x0E] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x0F] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x10] = color.RGBA{R: 152, G: 150, B: 152, A: 255}
	p.lookupPalette[0x11] = color.RGBA{R: 8, G: 76, B: 196, A: 255}
	p.lookupPalette[0x12] = color.RGBA{R: 48, G: 50, B: 236, A: 255}
	p.lookupPalette[0x13] = color.RGBA{R: 92, G: 30, B: 228, A: 255}
	p.lookupPalette[0x14] = color.RGBA{R: 136, G: 20, B: 176, A: 255}
	p.lookupPalette[0x15] = color.RGBA{R: 160, G: 20, B: 100, A: 255}
	p.lookupPalette[0x16] = color.RGBA{R: 152, G: 34, B: 32, A: 255}
	p.lookupPalette[0x17] = color.RGBA{R: 120, G: 60, B: 0, A: 255}
	p.lookupPalette[0x18] = color.RGBA{R: 84, G: 90, B: 0, A: 255}
	p.lookupPalette[0x19] = color.RGBA{R: 40, G: 114, B: 0, A: 255}
	p.lookupPalette[0x1A] = color.RGBA{R: 8, G: 124, B: 0, A: 255}
	p.lookupPalette[0x1B] = color.RGBA{R: 0, G: 118, B: 40, A: 255}
	p.lookupPalette[0x1C] = color.RGBA{R: 0, G: 102, B: 120, A: 255}
	p.lookupPalette[0x1D] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x1E] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x1F] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x20] = color.RGBA{R: 236, G: 238, B: 236, A: 255}
	p.lookupPalette[0x21] = color.RGBA{R: 76, G: 154, B: 236, A: 255}
	p.lookupPalette[0x22] = color.RGBA{R: 120, G: 124, B: 236, A: 255}
	p.lookupPalette[0x23] = color.RGBA{R: 176, G: 98, B: 236, A: 255}
	p.lookupPalette[0x24] = color.RGBA{R: 228, G: 84, B: 236, A: 255}
	p.lookupPalette[0x25] = color.RGBA{R: 236, G: 88, B: 180, A: 255}
	p.lookupPalette[0x26] = color.RGBA{R: 236, G: 106, B: 100, A: 255}
	p.lookupPalette[0x27] = color.RGBA{R: 212, G: 136, B: 32, A: 255}
	p.lookupPalette[0x28] = color.RGBA{R: 160, G: 170, B: 0, A: 255}
	p.lookupPalette[0x29] = color.RGBA{R: 116, G: 196, B: 0, A: 255}
	p.lookupPalette[0x2A] = color.RGBA{R: 76, G: 208, B: 32, A: 255}
	p.lookupPalette[0x2B] = color.RGBA{R: 56, G: 204, B: 108, A: 255}
	p.lookupPalette[0x2C] = color.RGBA{R: 56, G: 180, B: 204, A: 255}
	p.lookupPalette[0x2D] = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	p.lookupPalette[0x2E] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x2F] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x30] = color.RGBA{R: 236, G: 238, B: 236, A: 255}
	p.lookupPalette[0x31] = color.RGBA{R: 168, G: 204, B: 236, A: 255}
	p.lookupPalette[0x32] = color.RGBA{R: 188, G: 188, B: 236, A: 255}
	p.lookupPalette[0x33] = color.RGBA{R: 212, G: 178, B: 236, A: 255}
	p.lookupPalette[0x34] = color.RGBA{R: 236, G: 174, B: 236, A: 255}
	p.lookupPalette[0x35] = color.RGBA{R: 236, G: 174, B: 212, A: 255}
	p.lookupPalette[0x36] = color.RGBA{R: 236, G: 180, B: 176, A: 255}
	p.lookupPalette[0x37] = color.RGBA{R: 228, G: 196, B: 144, A: 255}
	p.lookupPalette[0x38] = color.RGBA{R: 204, G: 210, B: 120, A: 255}
	p.lookupPalette[0x39] = color.RGBA{R: 180, G: 222, B: 120, A: 255}
	p.lookupPalette[0x3A] = color.RGBA{R: 168, G: 226, B: 144, A: 255}
	p.lookupPalette[0x3B] = color.RGBA{R: 152, G: 226, B: 180, A: 255}
	p.lookupPalette[0x3C] = color.RGBA{R: 160, G: 214, B: 228, A: 255}
	p.lookupPalette[0x3D] = color.RGBA{R: 160, G: 162, B: 160, A: 255}
	p.lookupPalette[0x3E] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	p.lookupPalette[0x3F] = color.RGBA{R: 0, G: 0, B: 0, A: 255}
}
