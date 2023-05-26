package olc2c02

import (
	"nes_emulator/mlog"
)

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

func (p *Ppu) loadBgShifters() {
	attribLowConcat, attribHighConcat := 0, 0
	if p.bgNextTileAttrib&0b01 != 0 {
		attribLowConcat = 0xFF
	}
	if p.bgNextTileAttrib&0b10 != 0 {
		attribHighConcat = 0xFF
	}

	p.bgShifterPatternLow = (p.bgShifterPatternLow & 0xFF00) | uint16(p.bgNextTileLsb)
	p.bgShifterPatternHigh = (p.bgShifterPatternHigh & 0xFF00) | uint16(p.bgNextTileMsb)
	p.bgShifterAttribLow = (p.bgShifterAttribLow & 0xFF00) | uint16(attribLowConcat)
	p.bgShifterAttribHigh = (p.bgShifterAttribHigh & 0xFF00) | uint16(attribHighConcat)
}

func (p *Ppu) updateShifters() {
	if !p.regMask.GetFlag(regMaskRenderBG) {
		return
	}
	p.bgShifterPatternLow <<= 1
	p.bgShifterPatternHigh <<= 1
	p.bgShifterAttribLow <<= 1
	p.bgShifterAttribHigh <<= 1
}

func (p *Ppu) Clock() {
	// Clock timing table: https://www.nesdev.org/wiki/File:Ppu.svg
	// Behaviour detail: https://www.nesdev.org/wiki/PPU_rendering

	// This wrap all visible scanline operation
	if p.scanLine >= -1 && p.scanLine < 240 {
		// Check if we leave vertical blanking
		if p.scanLine == -1 && p.cycle == 1 {
			p.regStat.SetFlag(regStatVertZeroBlank, false)
		}

		// Code ref: https://github.com/OneLoneCoder/olcNES/blob/master/Part%20%234%20-%20PPU%20Backgrounds/olc2C02.cpp
		// Preload ppu with information to render next 8 tile pixel
		if (p.cycle >= 2 && p.cycle < 258) || (p.cycle >= 321 && p.cycle < 338) {
			// Update shifter every visible cycle
			p.updateShifters()

			switch (p.cycle - 1) % 8 { // 8 cycle wrap around
			case 0:
				p.loadBgShifters() // load next 8 pixel worth of data
				p.bgNextTileId = p.PRead(0x2000 | p.regLoopyVram.data&0x0FFF)
			case 2:
				p.bgNextTileAttrib = p.PRead(0x23C0 |
					(uint16(p.regLoopyVram.GetFlag(regLoopyNametableY)) << 11) |
					(uint16(p.regLoopyVram.GetFlag(regLoopyNametableX)) << 10) |
					((uint16(p.regLoopyVram.GetFlag(regLoopyCoarseY)) >> 2) << 3) |
					(uint16(p.regLoopyVram.GetFlag(regLoopyCoarseX)) >> 2))
				if p.regLoopyVram.GetFlag(regLoopyCoarseY)&0x02 != 0 {
					p.bgNextTileAttrib >>= 4
				}
				if p.regLoopyVram.GetFlag(regLoopyCoarseX)&0x02 != 0 {
					p.bgNextTileAttrib >>= 2
				}
				p.bgNextTileAttrib &= 0x03
			case 4:
				regCtrlPtrBgVal := uint16(0)
				if p.regCtrl.GetFlag(regCtrlPatternBG) {
					regCtrlPtrBgVal = 1 << 12
				}
				p.bgNextTileLsb = p.PRead(regCtrlPtrBgVal +
					(uint16(p.bgNextTileId) << 4) +
					uint16(p.regLoopyVram.GetFlag(regLoopyFineY)) +
					0)
			case 6: // exactly the same with offset 8
				regCtrlPtrBgVal := uint16(0)
				if p.regCtrl.GetFlag(regCtrlPatternBG) {
					regCtrlPtrBgVal = 1 << 12
				}
				p.bgNextTileMsb = p.PRead(regCtrlPtrBgVal +
					(uint16(p.bgNextTileId) << 4) +
					uint16(p.regLoopyVram.GetFlag(regLoopyFineY)) +
					8)
			case 7:
				// Must going on next tile
				p.regLoopyVram.IncrementScrollX(p.regMask)
			}
		}
		if p.cycle == 256 { // done with visible row, go to next y
			p.regLoopyVram.IncrementScrollY(p.regMask)
		}
		if p.cycle == 257 { // we incremented our y but x still wrong, so reset it
			p.regLoopyVram.TransferAddressX(p.regMask, p.regLoopyTram)
		}
		if p.scanLine == -1 && p.cycle >= 280 && p.cycle < 305 { // ready for new frame
			p.regLoopyVram.TransferAddressY(p.regMask, p.regLoopyTram)
		}
	}

	if p.scanLine == 240 {
		// Post rendering scanline, does nothing
	}

	// Check if we reach the first out of range scanline and emit NMI if enabled
	if p.scanLine == 241 && p.cycle == 1 {
		p.regStat.SetFlag(regStatVertZeroBlank, true)
		if p.regCtrl.GetFlag(regCtrlEnableNMI) {
			p.hasNmi = true
		}
	}

	// Render background
	if p.cycle < 256 && p.scanLine >= 0 && p.scanLine < 240 {
		if p.regMask.GetFlag(regMaskRenderBG) {
			bitMux := 0x8000 >> p.scrollFineX // select shift register
			bgPixel := uint8(0)
			if p.bgShifterPatternLow&uint16(bitMux) != 0 {
				bgPixel |= 0b01
			}
			if p.bgShifterPatternHigh&uint16(bitMux) != 0 {
				bgPixel |= 0b10
			}
			bgPalette := uint8(0)
			if p.bgShifterAttribLow&uint16(bitMux) != 0 {
				bgPalette |= 0b01
			}
			if p.bgShifterAttribHigh&uint16(bitMux) != 0 {
				bgPalette |= 0b10
			}
			pixelColor := p.getColorFromPaletteRam(bgPalette, bgPixel)
			i := 4 * (int(p.cycle) + int(p.scanLine)*256)
			p.screenDisplayBuf[i] = pixelColor.R
			p.screenDisplayBuf[i+1] = pixelColor.G
			p.screenDisplayBuf[i+2] = pixelColor.B
			p.screenDisplayBuf[i+3] = pixelColor.A
		}
	}

	//// Fake some noise for now
	//if p.cycle < 256 && p.scanLine >= 0 && p.scanLine < 240 {
	//	i := 4 * (int(p.cycle) + int(p.scanLine)*256)
	//	if rand.Float32() < 0.5 {
	//		p.screenDisplayBuf[i] = 255
	//		p.screenDisplayBuf[i+1] = 255
	//		p.screenDisplayBuf[i+2] = 255
	//	} else {
	//		p.screenDisplayBuf[i] = 0
	//		p.screenDisplayBuf[i+1] = 0
	//		p.screenDisplayBuf[i+2] = 0
	//	}
	//}

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

func (p *Ppu) Reset() {
	mlog.L.Info("Resetting ppu")
	p.scanLine = 0
	p.cycle = 0
	p.frameCompleted = false
	p.hasNmi = false
	p.regCtrl.data = 0
	p.regMask.data = 0
	p.regStat.data = 0
	p.regLoopyTram.data = 0
	p.regLoopyVram.data = 0
	p.addressLatch = false
	p.dataBuffer = 0
}
