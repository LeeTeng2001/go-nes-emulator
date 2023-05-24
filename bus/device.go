package bus

// CpuDevice has the cpu interface
// communication with main bus
type CpuDevice interface {
	CWrite(addr uint16, data uint8)
	CRead(addr uint16) (data uint8)
	Reset()
	Clock()
}

// PpuDevice has the ppu interface
// communication with ppu bus
type PpuDevice interface {
	CWrite(addr uint16, data uint8)
	CRead(addr uint16) (data uint8)
	PWrite(addr uint16, data uint8)
	PRead(addr uint16) (data uint8)
	Clock()
}
