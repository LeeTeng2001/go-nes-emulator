package olc2c02

import (
	"nes_emulator/mlog"
	"nes_emulator/utils"
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
	if p.regMask.GetFlag(regMaskRenderBG) {
		p.bgShifterPatternLow <<= 1
		p.bgShifterPatternHigh <<= 1
		p.bgShifterAttribLow <<= 1
		p.bgShifterAttribHigh <<= 1
	}
	if p.regMask.GetFlag(regMaskRenderSprite) && p.cycle >= 1 && p.cycle < 258 {
		// Only shift if coordinate == 0: ie in range for drawing
		for i := uint8(0); i < p.nextLineSpriteCount; i++ {
			if p.nextLineScanlineSprites[i].x > 0 {
				p.nextLineScanlineSprites[i].x--
			} else {
				p.nextLineSpriteShiftPtrnLo[i] <<= 1
				p.nextLineSpriteShiftPtrnHi[i] <<= 1
			}
		}
	}
}

func (p *Ppu) Clock() {
	// Clock timing table: https://www.nesdev.org/wiki/File:Ppu.svg
	// Behaviour detail: https://www.nesdev.org/wiki/PPU_rendering

	// This wrap all visible scanline operation
	if p.scanLine >= -1 && p.scanLine < 240 {
		// Background rendering ------------------------------------------------
		// Check if we leave vertical blanking, at the start of new frame
		if p.scanLine == -1 && p.cycle == 1 {
			p.regStat.SetFlag(regStatVertZeroBlank, false)
			p.regStat.SetFlag(regStatSpriteOverflow, false)
			for i := 0; i < 8; i++ {
				p.nextLineSpriteShiftPtrnLo[i] = 0
				p.nextLineSpriteShiftPtrnHi[i] = 0
			}
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
		if p.cycle == 338 || p.cycle == 340 { // last cycle, read next row tile
			p.bgNextTileId = p.PRead(0x2000 | p.regLoopyVram.data&0x0FFF)
		}

		// Foreground rendering ------------------------------------------------
		// This rendering doesn't follow the sprite evaluation in NES
		// we evaluate everything at specific cycle to make things simple
		if p.cycle == 257 && p.scanLine >= 0 { // end of scanline + first out of visible range
			// 1. clear next line buffer
			for idx, s := range p.nextLineScanlineSprites {
				s.x = 0xFF
				s.y = 0xFF
				s.id = 0xFF
				s.attribute = 0xFF
				p.nextLineSpriteShiftPtrnLo[idx] = 0
				p.nextLineSpriteShiftPtrnHi[idx] = 0
			}
			p.nextLineSpriteCount = 0

			// 2. Evaluate next line visible sprites by comparing y diff for all OAM sprites
			for nOAMEntry := 0; nOAMEntry < 64; nOAMEntry++ {
				yDiff := p.scanLine - int16(p.oamMem[nOAMEntry].y)
				spriteYSize := int16(8)
				if p.regCtrl.GetFlag(regCtrlSpriteSize) {
					spriteYSize = 16
				}

				if yDiff >= 0 && yDiff < spriteYSize { // scanline in range
					if p.nextLineSpriteCount < 8 { // next line has storage left
						p.nextLineScanlineSprites[p.nextLineSpriteCount] = p.oamMem[nOAMEntry]
						p.nextLineSpriteCount++
					} else { // sprite overflow
						p.regStat.SetFlag(regStatSpriteOverflow, true)
						mlog.L.Warn("Sprite overflow for single scanline")
					}
				}
			}
		}

		// Get data from pattern memory
		if p.cycle == 340 { // end of cycle
			// https://www.nesdev.org/wiki/PPU_OAM
			for i := uint8(0); i < p.nextLineSpriteCount; i++ {
				var patternBitsLow, patternBitsHigh uint8
				var patternAddrLow, patternAddrHigh uint16
				patternTableAddr := uint16(0)

				// Workflow: choose pattern table, choose tile, choose row

				if !p.regCtrl.GetFlag(regCtrlSpriteSize) {
					// 8x8 sprite
					// Pattern table (implied by control reg)
					if p.regCtrl.GetFlag(regCtrlPatternSprite) {
						patternTableAddr = 1 << 12
					}
					if p.nextLineScanlineSprites[i].attribute&0x80 == 0 {
						// normal
						patternAddrLow = patternTableAddr |
							(uint16(p.nextLineScanlineSprites[i].id) << 4) |
							(uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y))
					} else {
						// flipped vert
						patternAddrLow = patternTableAddr |
							(uint16(p.nextLineScanlineSprites[i].id) << 4) |
							(7 - (uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y)))
					}
				} else {
					// 8x16 sprite
					// Pattern table (implied by id)
					if p.nextLineScanlineSprites[i].id&1 != 0 {
						patternTableAddr = 1 << 12
					}
					if p.nextLineScanlineSprites[i].attribute&0x80 == 0 {
						// normal
						// work out top half or bottom half
						if p.scanLine-int16(p.nextLineScanlineSprites[i].y) < 8 {
							patternAddrLow = patternTableAddr |
								(uint16(p.nextLineScanlineSprites[i].id&0xFE) << 4) |
								((uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y)) & 0x07)
						} else {
							patternAddrLow = patternTableAddr |
								((uint16(p.nextLineScanlineSprites[i].id&0xFE) + 1) << 4) |
								((uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y)) & 0x07)
						}
					} else {
						// flipped vert
						// work out top half or bottom half
						if p.scanLine-int16(p.nextLineScanlineSprites[i].y) < 8 {
							patternAddrLow = patternTableAddr |
								((uint16(p.nextLineScanlineSprites[i].id&0xFE) + 1) << 4) |
								((7 - (uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y))) & 0x07)
						} else {
							patternAddrLow = patternTableAddr |
								(uint16(p.nextLineScanlineSprites[i].id&0xFE) << 4) |
								((7 - (uint16(p.scanLine) - uint16(p.nextLineScanlineSprites[i].y))) & 0x07)
						}
					}
				}

				// Get high and read corresponding data
				patternAddrHigh = patternAddrLow + 8
				patternBitsLow = p.PRead(patternAddrLow)
				patternBitsHigh = p.PRead(patternAddrHigh)

				// Determine horizontal flip
				if p.nextLineScanlineSprites[i].attribute&0x40 != 0 {
					patternBitsLow = utils.FlipByte(patternBitsLow)
					patternBitsHigh = utils.FlipByte(patternBitsHigh)
				}

				// Load into shift register
				p.nextLineSpriteShiftPtrnLo[i] = patternBitsLow
				p.nextLineSpriteShiftPtrnHi[i] = patternBitsHigh
			}
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

	// Write pixel to screen buffer
	if p.cycle < 256 && p.scanLine >= 0 && p.scanLine < 240 {
		// Background color and palette
		bgPixel := uint8(0)
		bgPalette := uint8(0)
		if p.regMask.GetFlag(regMaskRenderBG) {
			bitMux := 0x8000 >> p.scrollFineX // select shift register
			if p.bgShifterPatternLow&uint16(bitMux) != 0 {
				bgPixel |= 0b01
			}
			if p.bgShifterPatternHigh&uint16(bitMux) != 0 {
				bgPixel |= 0b10
			}
			if p.bgShifterAttribLow&uint16(bitMux) != 0 {
				bgPalette |= 0b01
			}
			if p.bgShifterAttribHigh&uint16(bitMux) != 0 {
				bgPalette |= 0b10
			}
		}

		// Foreground color and palette
		fgPixel := uint8(0)
		fgPalette := uint8(0)
		fgPriority := uint8(0)
		if p.regMask.GetFlag(regMaskRenderSprite) {
			// Find the sprite from the highest render priority
			for i := uint8(0); i < p.nextLineSpriteCount; i++ {
				if p.nextLineScanlineSprites[i].x == 0 {
					if p.nextLineSpriteShiftPtrnLo[i]&0x80 != 0 {
						fgPixel |= 0b01
					}
					if p.nextLineSpriteShiftPtrnHi[i]&0x80 != 0 {
						fgPixel |= 0b10
					}
					// Unlike bg, sprite palette is contains inside attribute instead of region
					// Offset by 4 because first 4 is reserved for background
					fgPalette = p.nextLineScanlineSprites[i].attribute&0x03 + 0x04
					if p.nextLineScanlineSprites[i].attribute&0x20 == 0 { // allow sprite to go behind bg
						fgPriority = 1
					}

					// found top priority sprite with non-transparent pixel
					if fgPixel != 0 {
						break
					}
				}
			}
		}

		// 0, 0 is the default id both fg and bg is transparent
		finalPixel := uint8(0)
		finalPalette := uint8(0)

		// competition to determine pixel output
		if bgPixel == 0 && fgPixel != 0 {
			finalPixel = fgPixel
			finalPalette = fgPalette
		} else if bgPixel != 0 && fgPixel == 0 {
			finalPixel = bgPixel
			finalPalette = bgPalette
		} else if bgPixel != 0 && fgPixel != 0 { // both visible, check priority
			if fgPriority != 0 {
				finalPixel = fgPixel
				finalPalette = fgPalette
			} else {
				finalPixel = bgPixel
				finalPalette = bgPalette
			}
		}

		pixelColor := p.getColorFromPaletteRam(finalPalette, finalPixel)
		i := 4 * (int(p.cycle) + int(p.scanLine)*256)
		p.screenDisplayBuf[i] = pixelColor.R
		p.screenDisplayBuf[i+1] = pixelColor.G
		p.screenDisplayBuf[i+2] = pixelColor.B
		p.screenDisplayBuf[i+3] = pixelColor.A
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
