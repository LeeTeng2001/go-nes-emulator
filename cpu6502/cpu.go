package cpu6502

import "nes_emulator/bus"

// Compile time interface check
var _ bus.Device = (*Cpu)(nil)

const (
	irqHandlerAddr     = 0xFFFE
	nmiHandlerAddr     = 0xFFFA
	baseStackOffset    = 0x0100
	resetStackPtrAddr  = 0xFD
	resetPCAddr        = 0xFFFC
	resetRequiredCycle = 8
)

type Cpu struct {
	b *bus.Bus
	// Registers
	regA      uint8
	regX      uint8
	regY      uint8
	regStatus uint8
	regStkPtr uint8
	regPC     uint16
	// internal
	fetchedData uint8
	addrAbs     uint16
	addrRel     uint16
	opcode      uint8
	cyclesLeft  uint8
	// lookup table
	insLookup []instruction
}

func New() *Cpu {
	newCpu := &Cpu{}
	newCpu.initInstLookup()
	return newCpu
}

func (c *Cpu) Write(addr uint16, data uint8) {
	c.b.Write(addr, data)
}

func (c *Cpu) Read(addr uint16) (data uint8) {
	return c.b.Read(addr)
}

func (c *Cpu) ConnectBus(b *bus.Bus) {
	c.b = b
}

func (c *Cpu) getStatus(flag uint8) bool {
	return c.regStatus&flag > 0
}

func (c *Cpu) setStatus(flag uint8, setBit bool) {
	if setBit {
		c.regStatus |= flag
	} else {
		c.regStatus &= ^flag
	}
}
