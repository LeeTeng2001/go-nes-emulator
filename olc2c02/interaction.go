package olc2c02

import "nes_emulator/mlog"

func (p *Ppu) CWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: // Mask
		break
	case 0x0002: // Status
		break
	case 0x0003: // OAM Address
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		break
	default:
		mlog.L.Fatalf("Invalid control code %d encountered at ppu!", addr)
	}
}

func (p *Ppu) CRead(addr uint16) (data uint8) {
	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: // Mask
		break
	case 0x0002: // Status
		break
	case 0x0003: // OAM Address
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		break
	default:
		mlog.L.Fatalf("Invalid control code %d encountered at ppu!", addr)
	}
	return data
}

func (p *Ppu) PWrite(addr uint16, data uint8) {
	addr &= 0x3FFF // Hardware limit
	if p.disk.PWrite(addr, data) {
		// Empty for now
	}
}

func (p *Ppu) PRead(addr uint16) (data uint8) {
	addr &= 0x3FFF // Hardware limit
	if p.disk.PRead(addr, &data) {
		// Empty for now
		return data
	}
	return 0
}
