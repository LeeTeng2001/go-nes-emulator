package olc2c02

import (
	"math/rand"
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

func (p *Ppu) Reset() {
	mlog.L.Info("Resetting ppu")
	p.scanLine = 0
	p.cycle = 0
	p.frameCompleted = false
	p.hasNmi = false
	p.regCtrl.data = 0
	p.regMask.data = 0
	p.regStat.data = 0
	p.hasWriteBuffer = false
	p.addrCombined = 0
	p.dataBuffer = 0
}
