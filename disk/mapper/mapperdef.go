package mapper

// IMapper Takes request address and map it to the rom like we've read it from the file
type IMapper interface {
	CpuMapRead(addr uint16, mappedAddr *uint32) bool
	CpuMapWrite(addr uint16, mappedAddr *uint32) bool
	PpuMapRead(addr uint16, mappedAddr *uint32) bool
	PpuMapWrite(addr uint16, mappedAddr *uint32) bool
}

type mapper struct {
	prgBanks uint8
	chrBanks uint8
}
