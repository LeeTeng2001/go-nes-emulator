package olc2c02

import "nes_emulator/bus"

// Compile time interface check
var _ bus.BothDevice = (*Ppu)(nil)

type Ppu struct {
	b *bus.Bus
	// Ram
	tableName    [2][1024]uint8
	tablePattern [32]uint8
}

func New() *Ppu {
	return &Ppu{}
}

func (p Ppu) ConnectBus(b *bus.Bus) {
	p.b = b
}

func (p Ppu) CWrite(addr uint16, data uint8) {
	//TODO implement me
	panic("implement me")
}

func (p Ppu) CRead(addr uint16) (data uint8) {
	//TODO implement me
	panic("implement me")
}

func (p Ppu) PWrite(addr uint16, data uint8) {
	//TODO implement me
	panic("implement me")
}

func (p Ppu) PRead(addr uint16) (data uint8) {
	//TODO implement me
	panic("implement me")
}
