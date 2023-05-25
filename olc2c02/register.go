package olc2c02

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

type register struct {
	data uint8
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
