package mapper

// Definition: https://www.nesdev.org/wiki/INES_Mapper_000
// All banks are fixed, very simple

// Compile time interface check
var _ IMapper = (*M0)(nil)

type M0 struct {
	mapper
}

func (m *M0) CpuMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 { // to the end of address
		if m.prgBanks > 1 { // more than 1 bank?
			addr = addr & 0x7FFF
		} else { // mirror bank according to specification
			addr = addr & 0x3FFF
		}
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m *M0) CpuMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 {
		// Same as read
		if m.prgBanks > 1 {
			addr = addr & 0x7FFF
		} else {
			addr = addr & 0x3FFF
		}
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m *M0) PpuMapRead(addr uint16, mappedAddr *uint32) bool {
	// No mapping for ppu, directly map to chr rom
	if addr <= 0x1FFF {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m *M0) PpuMapWrite(addr uint16, mappedAddr *uint32) bool {
	// doesn't make sense, but if we have 0 character banks treat it as ram
	if addr <= 0x1FFF && m.chrBanks == 0 {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func NewM0(prgBanks uint8, chrBanks uint8) *M0 {
	m := &M0{}
	m.prgBanks = prgBanks
	m.chrBanks = chrBanks
	return m
}
