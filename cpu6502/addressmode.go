package cpu6502

// Reference: https://www.nesdev.org/wiki/CPU_addressing_modes
// Code ref: https://github.com/OneLoneCoder/olcNES/blob/master/Part%232%20-%20CPU/olc6502.cpp

func (c *Cpu) IMP() uint8 {
	// Implied
	c.fetchedData = c.regA
	return 0
}

func (c *Cpu) IMM() uint8 {
	// Immediate
	c.addrAbs = c.regPC
	c.regPC++
	return 0
}

func (c *Cpu) ZP0() uint8 {
	// address can be split into 0xAABB, where AA is page and BB is offset
	// Zero page addressing means page = 0
	c.addrAbs = uint16(c.read(c.regPC))
	c.regPC++
	return 0
}

func (c *Cpu) ZPX() uint8 {
	// Useful for iterating array
	c.addrAbs = uint16(c.read(c.regPC) + c.regX)
	c.regPC++
	return 0
}

func (c *Cpu) ZPY() uint8 {
	c.addrAbs = uint16(c.read(c.regPC) + c.regY)
	c.regPC++
	return 0
}

func (c *Cpu) REL() uint8 {
	// Only apply for branching instruction, can't jump further than 127
	c.addrRel = uint16(c.read(c.regPC))
	c.regPC++
	// check jumping backward
	if c.addrRel&0x80 != 0 {
		c.addrRel = c.addrRel | 0xFF00
	}
	return 0
}

func (c *Cpu) ABS() uint8 {
	low := c.read(c.regPC)
	c.regPC++
	high := c.read(c.regPC)
	c.regPC++
	c.addrAbs = (uint16(high) << 8) | uint16(low)
	return 0
}

func (c *Cpu) ABX() uint8 {
	low := c.read(c.regPC)
	c.regPC++
	high := c.read(c.regPC)
	c.regPC++
	c.addrAbs = (uint16(high) << 8) | uint16(low)
	c.addrAbs += uint16(c.regX)

	// If address have move to another page after addition we need extra clock cycle
	if c.addrAbs&0xFF00 != (uint16(high) << 8) {
		return 1
	}
	return 0
}

func (c *Cpu) ABY() uint8 {
	low := c.read(c.regPC)
	c.regPC++
	high := c.read(c.regPC)
	c.regPC++
	c.addrAbs = (uint16(high) << 8) | uint16(low)
	c.addrAbs += uint16(c.regY)

	// If address have move to another page after addition we need extra clock cycle
	if c.addrAbs&0xFF00 != (uint16(high) << 8) {
		return 1
	}
	return 0

}

func (c *Cpu) IND() uint8 {
	// 6502 way of implementing pointer
	ptrLow := c.read(c.regPC)
	c.regPC++
	ptrHigh := c.read(c.regPC)
	c.regPC++
	ptr := (uint16(ptrHigh) << 8) | uint16(ptrLow)

	// Overflow if ptrLow == 0xFF, simulate page boundary
	if ptrLow == 0xFF {
		c.addrAbs = (uint16(c.read(ptr&0xFF00)) << 8) | uint16(c.read(ptr))
	} else {
		c.addrAbs = (uint16(c.read(ptr+1)) << 8) | uint16(c.read(ptr))
	}

	return 0
}

func (c *Cpu) IZX() uint8 {
	// Indirect access with zero page
	t := c.read(c.regPC)
	c.regPC++
	low := c.read(uint16(t+c.regX) & 0x00FF)
	high := c.read(uint16(t+c.regX+1) & 0x00FF)
	c.addrAbs = (uint16(high) << 8) | uint16(low)
	return 0
}

func (c *Cpu) IZY() uint8 {
	// Indirect access with zero page
	t := c.read(c.regPC)
	c.regPC++
	low := c.read(uint16(t) & 0x00FF)
	high := c.read(uint16(t+1) & 0x00FF)
	c.addrAbs = (uint16(high) << 8) | uint16(low)
	c.addrAbs += uint16(c.regY)

	// If address have move to another page after addition we need extra clock cycle
	if c.addrAbs&0xFF00 != (uint16(high) << 8) {
		return 1
	}
	return 0
}
