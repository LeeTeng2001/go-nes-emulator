package cpu6502

import (
	"github.com/charmbracelet/log"
	"nes_emulator/bus"
	"nes_emulator/mlog"
	"os"
)

// Compile time interface check
var _ bus.CpuDevice = (*Cpu)(nil)

const (
	irqHandlerAddr     = 0xFFFE
	nmiHandlerAddr     = 0xFFFA
	baseStackOffset    = 0x0100
	resetStackPtrAddr  = 0xFD
	resetPCAddr        = 0xFFFC
	resetRequiredCycle = 8
)

const (
	cpuLogFile = "cpu.log"
)

type Cpu struct {
	// Our ds
	b          *bus.Bus
	cpuFileLog *log.Logger
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
	mlog.L.Info("Cpu is initialised")
	return newCpu
}

func NewDebug() *Cpu {
	newCpu := &Cpu{}
	newCpu.initInstLookup()

	file, err := os.Create(cpuLogFile)
	if err != nil {
		log.Fatalf("cannot open %v", cpuLogFile)
	}
	newCpu.cpuFileLog = log.NewWithOptions(file, log.Options{
		ReportTimestamp: true,
	})

	mlog.L.Info("Cpu is initialised with debug enabled")
	return newCpu
}

func (c *Cpu) ConnectBus(b *bus.Bus) {
	c.b = b
}

// Helper functions ------------------------------------------

func (c *Cpu) write(addr uint16, data uint8) {
	c.b.CWrite(addr, data)
}

func (c *Cpu) read(addr uint16) (data uint8) {
	return c.b.CRead(addr)
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
