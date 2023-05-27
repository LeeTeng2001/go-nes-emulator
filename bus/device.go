package bus

import "nes_emulator/disk"

// CpuDevice has the cpu interface
// communication with main bus
type CpuDevice interface {
	Reset()
	Clock()
	Nmi()
}

type ApuDevice interface {
	CWrite(addr uint16, data uint8)
	CRead(addr uint16) (data uint8)
	Reset()
	Clock()
	GetOutputSample() float64
}

// PpuDevice has the ppu interface
// communication with ppu bus
type PpuDevice interface {
	CWrite(addr uint16, data uint8)
	CRead(addr uint16) (data uint8)
	PWrite(addr uint16, data uint8)
	PRead(addr uint16) (data uint8)
	ConnectDisk(nesDisk *disk.NesDisk)
	Reset()
	Clock()
	CheckNmiAndTurnOff() bool
	SetOAMAsBytes(addr uint8, data uint8) // for DMA
}
