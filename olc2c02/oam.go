package olc2c02

// Stored as nes expected structure
type spriteObjEntry struct {
	y         uint8 // y pos of sprite
	id        uint8 // tile id from pattern memory
	attribute uint8 // how should sprite be rendered?
	x         uint8 // x pos of sprite
}

// GetOAMAsBytes access sprite obj array like a sequence of byte array
func (p *Ppu) GetOAMAsBytes(addr uint8) (data uint8) {
	arrLoc := addr / 4
	offset := addr % 4
	switch offset {
	case 0:
		return p.oamMem[arrLoc].y
	case 1:
		return p.oamMem[arrLoc].id
	case 2:
		return p.oamMem[arrLoc].attribute
	case 3:
		return p.oamMem[arrLoc].x
	}
	return 0
}

func (p *Ppu) SetOAMAsBytes(addr uint8, data uint8) {
	arrLoc := addr / 4
	offset := addr % 4
	switch offset {
	case 0:
		p.oamMem[arrLoc].y = data
	case 1:
		p.oamMem[arrLoc].id = data
	case 2:
		p.oamMem[arrLoc].attribute = data
	case 3:
		p.oamMem[arrLoc].x = data
	}
}
