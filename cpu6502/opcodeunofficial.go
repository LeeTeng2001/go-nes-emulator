package cpu6502

// Unofficial opcode! http://www.ffd2.com/fridge/docs/6502-NMOS.extra.opcodes

func (c *Cpu) LAX() uint8 {
	// https://www.nesdev.org/wiki/Programming_with_unofficial_opcodes
	// Shortcut for LDA than TAX
	c.fetch()
	c.regA = c.fetchedData
	c.regX = c.regA
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 0
}

// TODO: Maybe implement other unofficial opcode
