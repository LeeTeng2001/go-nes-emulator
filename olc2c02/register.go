package olc2c02

import "nes_emulator/mlog"

// Definition ref: https://www.nesdev.org/wiki/PPU_registers

const (
	regCtrlNameTableX = 1 << iota
	regCtrlNameTableY
	regCtrlIncrementMode
	regCtrlPatternSprite
	regCtrlPatternBG
	regCtrlSpriteSize
	regCtrlSlaveMode
	regCtrlEnableNMI
)
const (
	regMaskGreyscale = 1 << iota
	regMaskRenderBGLeft
	regMaskRenderSpriteLeft
	regMaskRenderBG
	regMaskRenderSprite
	regMaskEnhanceRed
	regMaskEnhanceBlue
	regMaskEnhanceGreen
)
const (
	regStatSpriteOverflow = 1 << 5
	regStatSpriteZeroHit  = 1 << 6
	regStatVertZeroBlank  = 1 << 7
)

type loopyBitFields uint8

const (
	regLoopyCoarseX loopyBitFields = iota
	regLoopyCoarseY
	regLoopyNametableX
	regLoopyNametableY
	regLoopyFineY
)

type register struct {
	data uint8
}

type registerLoopy struct {
	data uint16
}

func (r *register) GetFlag(mask uint8) bool {
	if r.data&mask != 0 {
		return true
	}
	return false
}

func (r *register) SetFlag(mask uint8, on bool) {
	if on {
		r.data = mask | r.data
	} else {
		r.data = ^mask & r.data
	}
}

func (r *registerLoopy) GetFlag(mask loopyBitFields) uint8 {
	switch mask {
	case regLoopyCoarseX:
		return uint8(r.data & 0b11111)
	case regLoopyCoarseY:
		return uint8((r.data >> 5) & 0b11111)
	case regLoopyNametableX:
		return uint8((r.data >> 10) & 1)
	case regLoopyNametableY:
		return uint8((r.data >> 11) & 1)
	case regLoopyFineY:
		return uint8((r.data >> 12) & 0b111)
	default:
		mlog.L.Fatalf("Invalid mask code for loopy field: %x", mask)
	}
	return 0
}

func (r *registerLoopy) SetFlag(mask loopyBitFields, data uint8) {
	switch mask {
	case regLoopyCoarseX:
		r.data = ((r.data >> 5) << 5) | (uint16(data) & 0b11111)
	case regLoopyCoarseY:
		tmp := r.data & 0b11111
		r.data = ((r.data >> 10) << 10) | tmp | ((uint16(data) & 0b11111) << 5)
	case regLoopyNametableX:
		tmp := r.data & 0b1111111111
		r.data = ((r.data >> 11) << 11) | tmp | ((uint16(data) & 1) << 10)
	case regLoopyNametableY:
		tmp := r.data & 0b11111111111
		r.data = ((r.data >> 12) << 12) | tmp | ((uint16(data) & 1) << 11)
	case regLoopyFineY:
		tmp := r.data & 0b111111111111
		r.data = ((r.data >> 15) << 15) | tmp | ((uint16(data) & 0b111) << 12)
	default:
		mlog.L.Fatalf("Invalid mask code for loopy field: %x", mask)
	}
}

// Additional functions done on loopy register at the end of cycle group ------------------------

func (r *registerLoopy) IncrementScrollX(regMask register) {
	// No render, stop
	if !regMask.GetFlag(regMaskRenderBG) && !regMask.GetFlag(regMaskRenderSprite) {
		return
	}

	// A single name table is 32x30 tiles. As we increment horizontally
	// we may cross into a neighbouring nametable, or wrap around to
	// a neighbouring nametable
	coarseXVal := r.GetFlag(regLoopyCoarseX)
	if coarseXVal == 31 { // leaving table, wrap around
		r.SetFlag(regLoopyCoarseX, 0)
		tmp := r.GetFlag(regLoopyNametableX)
		r.SetFlag(regLoopyNametableX, ^tmp)
	} else { // inside current nametable
		r.SetFlag(regLoopyCoarseX, coarseXVal+1)
	}
}

func (r *registerLoopy) IncrementScrollY(regMask register) {
	// No render, stop
	if !regMask.GetFlag(regMaskRenderBG) && !regMask.GetFlag(regMaskRenderSprite) {
		return
	}

	// Incremnt y is complicated because last two rows of nametable is used as attribute tables
	fineYVal := r.GetFlag(regLoopyFineY)
	if fineYVal < 7 { // can increment inside a tile
		r.SetFlag(regLoopyFineY, fineYVal+1)
	} else { // Potential wrapping into other nametable
		// The coarse y offset is used to identify which
		// row of the nametable we want, and the fine
		// y offset is the specific "scanline"
		r.SetFlag(regLoopyFineY, 0) // reset
		coarseYVal := r.GetFlag(regLoopyCoarseY)
		if coarseYVal == 29 { // last vertical row, swap nametable
			r.SetFlag(regLoopyCoarseY, 0)
			r.SetFlag(regLoopyNametableY, ^r.GetFlag(regLoopyNametableY)) // flip
		} else if coarseYVal == 31 { // inside attribute memory, wrap around current nametable
			r.SetFlag(regLoopyCoarseY, 0)
		} else { // no special case
			r.SetFlag(regLoopyCoarseY, coarseYVal+1)
		}
	}
}

func (r *registerLoopy) TransferAddressX(regMask register, regLoopyTmp registerLoopy) {
	// No render, stop
	if !regMask.GetFlag(regMaskRenderBG) && !regMask.GetFlag(regMaskRenderSprite) {
		return
	}

	// Transfer the temporarily stored horizontal nametable access information
	// into the "pointer". Note that fine x scrolling is not part of the mechanism
	r.SetFlag(regLoopyNametableX, regLoopyTmp.GetFlag(regLoopyNametableX))
	r.SetFlag(regLoopyCoarseX, regLoopyTmp.GetFlag(regLoopyCoarseX))
}

func (r *registerLoopy) TransferAddressY(regMask register, regLoopyTmp registerLoopy) {
	// No render, stop
	if !regMask.GetFlag(regMaskRenderBG) && !regMask.GetFlag(regMaskRenderSprite) {
		return
	}

	// Transfer the temporarily stored horizontal nametable access information
	// into the "pointer". Note that fine x scrolling is not part of the mechanism
	r.SetFlag(regLoopyNametableY, regLoopyTmp.GetFlag(regLoopyNametableY))
	r.SetFlag(regLoopyCoarseY, regLoopyTmp.GetFlag(regLoopyCoarseY))
	r.SetFlag(regLoopyFineY, regLoopyTmp.GetFlag(regLoopyFineY))
}
