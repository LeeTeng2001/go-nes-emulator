package olc2c02

import "nes_emulator/mlog"

func (p *Ppu) CWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		p.regCtrl.data = data
	case 0x0001: // Mask
		p.regMask.data = data
	case 0x0002: // Status (can't write to)
		mlog.L.Fatal("You can't write to status register in ppu!")
	case 0x0003: // OAM Address
	case 0x0004: // OAM Data
	case 0x0005: // Scroll
	case 0x0006: // PPU Address (2 cycle to write full address)
		if !p.hasWriteBuffer {
			p.hasWriteBuffer = true
			p.addrCombined = (p.addrCombined & 0x00FF) | (uint16(data) << 8)
		} else {
			p.hasWriteBuffer = false
			p.addrCombined = (p.addrCombined & 0xFF00) | uint16(data)
		}
		break
	case 0x0007: // PPU Data (will auto increment addr to avoid tedious set and write on successive location)
		p.PWrite(p.addrCombined, data)
		// To speed up write, control register contains info about the orientation of auto increment address
		if p.regCtrl.GetFlag(regCtrlIncrementMode) {
			p.addrCombined += 32 // going down row
		} else {
			p.addrCombined++
		}
	default:
		mlog.L.Fatalf("Invalid control code %d encountered at ppu!", addr)
	}
}

func (p *Ppu) CRead(addr uint16) (data uint8) {
	switch addr {
	case 0x0000: // Control
		return p.regCtrl.data
	case 0x0001: // Mask
		return p.regMask.data
	case 0x0002: // Status (act of reading will change the state of device!)
		// interested in top 3 bits (last 5 bits usually noise or last operation leftover)
		//p.regStat.SetFlag(regStatVertZeroBlank, true)
		data = (p.regStat.data & 0xE0) | (p.dataBuffer & 0x1F)
		p.regStat.SetFlag(regStatVertZeroBlank, false)
		p.hasWriteBuffer = false
		return data
	case 0x0003: // OAM Address
	case 0x0004: // OAM Data
	case 0x0005: // Scroll
	case 0x0006: // PPU Address
		mlog.L.Fatal("You can't read from address register in ppu!")
	case 0x0007: // PPU Data (2 cycle for most cases except for palette)
		data = p.dataBuffer
		p.dataBuffer = p.PRead(p.addrCombined)
		if p.addrCombined >= PaletteRamOffset { // immediate output for palette
			data = p.dataBuffer
		}
		p.addrCombined++
		return data
	default:
		mlog.L.Fatalf("Invalid control code %d encountered at ppu!", addr)
	}
	return data
}

func (p *Ppu) PWrite(addr uint16, data uint8) {
	addr &= 0x3FFF // Hardware limit
	if p.disk.PWrite(addr, data) {
		// Empty for now
	} else if addr < 0x2000 { // pattern memory
		mlog.L.Warn("Pattern memory is usually read only but detected write")
		p.tablePatterns[(addr&0x1000)>>12][addr&0x0FFF] = data
	} else if addr >= 0x2000 && addr < 0x3F00 { // nametable memory, has mirroring capability
		addr = addr & 0x0FFF
		if !p.disk.MirrorHorizontal { // Vertical
			if addr <= 0x03FF {
				p.tableName[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				p.tableName[1][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				p.tableName[0][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				p.tableName[1][addr&0x03FF] = data
			}
		} else { // horizontal
			if addr <= 0x03FF {
				p.tableName[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				p.tableName[0][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				p.tableName[1][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				p.tableName[1][addr&0x03FF] = data
			}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // palette memory
		addr &= 0x001F
		switch addr {
		case 0x0010:
			addr = 0x0000
		case 0x0014:
			addr = 0x0004
		case 0x0018:
			addr = 0x0008
		case 0x001C:
			addr = 0x000C
		}
		p.tablePalette[addr] = data
		//mlog.L.Infof("Write to palette %X, %X", addr, data)
	}
}

func (p *Ppu) PRead(addr uint16) (data uint8) {
	addr &= 0x3FFF                 // Addressable range
	if p.disk.PRead(addr, &data) { // check mapper
		return data
	} else if addr < 0x2000 { // pattern memory, table id -> offset, but usually handle by disk
		return p.tablePatterns[(addr&0x1000)>>12][addr&0x0FFF]
	} else if addr >= 0x2000 && addr < 0x3F00 { // nametable memory, has mirroring capability
		// Specification: https://www.nesdev.org/wiki/PPU_nametables
		addr = addr & 0x0FFF
		if !p.disk.MirrorHorizontal { // Vertical
			if addr <= 0x03FF {
				data = p.tableName[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = p.tableName[1][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = p.tableName[0][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = p.tableName[1][addr&0x03FF]
			}
		} else { // horizontal
			if addr <= 0x03FF {
				data = p.tableName[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = p.tableName[0][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = p.tableName[1][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = p.tableName[1][addr&0x03FF]
			}
		}
		return data
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // palette memory
		// Get the offset
		addr &= 0x001F
		// Hardcode mirroring for specific address: https://www.nesdev.org/wiki/PPU_palettes
		switch addr {
		case 0x0010:
			addr = 0x0000
		case 0x0014:
			addr = 0x0004
		case 0x0018:
			addr = 0x0008
		case 0x001C:
			addr = 0x000C
		}
		// Read directly from predefined palette
		// TODO:  & (mask.grayscale ? 0x30 : 0x3F);
		return p.tablePalette[addr]
	}
	return 0
}
