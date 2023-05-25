package cpu6502

import (
	"nes_emulator/mlog"
	"reflect"
)

func (c *Cpu) Clock() {
	// read new instruction and update remaining cycle
	if c.cyclesLeft == 0 {
		// Log debug
		if c.cpuFileLog != nil {
			dis, _ := c.disassembleMemAtAddr(c.regPC)
			c.cpuFileLog.Printf("   %s %s", dis, c.getRegStateStr())
		}

		c.opcode = c.read(c.regPC)
		c.regPC++
		c.cyclesLeft = c.insLookup[c.opcode].cyclesRequired

		// Some instruction takes an additional cycle
		additionalCycle := c.insLookup[c.opcode].addrMode()
		additionalCycle = additionalCycle & c.insLookup[c.opcode].operate()
		c.cyclesLeft += additionalCycle
	}
	c.cyclesLeft--
}

func (c *Cpu) fetch() {
	// Get data if our opcode has value
	if reflect.ValueOf(c.insLookup[c.opcode].addrMode).Pointer() != reflect.ValueOf(c.IMP).Pointer() {
		c.fetchedData = c.read(c.addrAbs)
	}
}

// Reset cpu to a known state
func (c *Cpu) Reset() {
	mlog.L.Info("Resetting cpu")
	c.regA = 0
	c.regX = 0
	c.regY = 0
	c.regStkPtr = resetStackPtrAddr
	c.regStatus = 0 | FlagU | FlagI // TODO: Should set interrupt flag?
	c.addrAbs = 0
	c.addrRel = 0
	c.fetchedData = 0

	// Reset pc, the data at the location can be set by programmer
	low := c.read(resetPCAddr)
	high := c.read(resetPCAddr + 1)
	c.regPC = (uint16(high) << 8) | uint16(low)

	// A hardcoded reset cycle because it takes time
	c.cyclesLeft = resetRequiredCycle
}

// Irq interrupt request, maskable
func (c *Cpu) Irq() {
	if !c.getStatus(FlagI) {
		// Save pc, status reg at stack
		c.write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC>>8))
		c.regStkPtr--
		c.write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC))
		c.regStkPtr--

		c.setStatus(FlagB, false)
		c.setStatus(FlagU, true)
		c.setStatus(FlagI, true) // interrupt has occurred
		c.write(baseStackOffset+uint16(c.regStkPtr), c.regStatus)
		c.regStkPtr--

		// Service interrupt
		c.addrAbs = irqHandlerAddr
		low := c.read(c.addrAbs)
		high := c.read(c.addrAbs + 1)
		c.regPC = (uint16(high) << 8) | uint16(low)

		// Interrupt takes time
		c.cyclesLeft = resetRequiredCycle
	}
}

// Nmi non maskable interrupt, similar to Irq but the jump address is different
func (c *Cpu) Nmi() {
	// Save pc, status reg at stack
	c.write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC>>8))
	c.regStkPtr--
	c.write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC))
	c.regStkPtr--

	c.setStatus(FlagB, false)
	c.setStatus(FlagU, true)
	c.setStatus(FlagI, true) // interrupt has occurred
	c.write(baseStackOffset+uint16(c.regStkPtr), c.regStatus)
	c.regStkPtr--

	// Service interrupt, note that this address is different from irq
	c.addrAbs = nmiHandlerAddr
	low := c.read(c.addrAbs)
	high := c.read(c.addrAbs + 1)
	c.regPC = (uint16(high) << 8) | uint16(low)

	// Interrupt takes time
	c.cyclesLeft = resetRequiredCycle
}
