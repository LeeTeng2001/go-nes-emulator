package mapper

// Definition: https://www.nesdev.org/wiki/UxROM

// Compile time interface check
var _ IMapper = (*M2)(nil)

type M2 struct {
	mapper
	// Mapper registers state
	nPRGBankSelectLow  uint8
	nPRGBankSelectHigh uint8
}

func (m *M2) CpuMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xBFFF {
		*mappedAddr = uint32(m.nPRGBankSelectLow)*0x4000 + (uint32(addr) & 0x3FFF)
		return true
	} else if addr >= 0xC000 { // to the end of address, fixed to last bank
		*mappedAddr = uint32(m.nPRGBankSelectHigh)*0x4000 + (uint32(addr) & 0x3FFF)
		return true
	}
	return false
}

func (m *M2) CpuMapWrite(addr uint16, mappedAddr *uint32, data uint8) bool {
	if addr >= 0x8000 {
		// Note we handle write but we don't update ROM
		m.nPRGBankSelectLow = data & 0x0F
	}
	return false
}

func (m *M2) PpuMapRead(addr uint16, mappedAddr *uint32) bool {
	// No mapping for ppu, directly map to chr rom
	if addr < 0x2000 {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m *M2) PpuMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr < 0x2000 && m.chrBanks == 0 { // treating as ram
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func NewM2(prgBanks uint8, chrBanks uint8) *M2 {
	m := &M2{}
	m.prgBanks = prgBanks
	m.chrBanks = chrBanks
	return m
}

func (m *M2) Reset() {
	m.nPRGBankSelectLow = 0
	m.nPRGBankSelectHigh = m.prgBanks - 1
}
