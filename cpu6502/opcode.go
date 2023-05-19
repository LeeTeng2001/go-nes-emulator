package cpu6502

import (
	"nes_emulator/mlog"
	"reflect"
)

// Great reference: https://www.nesdev.org/obelisk-6502-guide/instructions.html

// arithmetic: fetch data, perform operation, update status reg, check potentially require additional clocks
// branch: will modify cc directly, if cross page an extra clock is required
// clear flag: simple clear flag instruction
// Stack instruction: push and pop, grow from a hardcode location
// return from interrupt: RTI, restore state
// break is similar to interrupt request
// compare instruction and set flag
// Increment, decrement at memory/reg
// Jump

func (c *Cpu) ADC() uint8 {
	// Addition, need to take care of overflow and negative
	// overflow formula derived from 49min mark: https://www.youtube.com/watch?v=8XmxKPJDGU0&list=PLrOv9FMX8xJHqMvSGB_9G9nZZ_4IgteYf&index=4
	c.fetch()
	tmpRes := uint16(c.regA) + uint16(c.fetchedData)
	if c.getStatus(FlagC) {
		tmpRes++
	}

	// Check flag
	c.setStatus(FlagC, tmpRes > 0xFF)
	c.setStatus(FlagZ, (tmpRes&0xFF) == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)
	c.setStatus(FlagV,
		(^uint16(c.regA^c.fetchedData)&(uint16(c.regA)^tmpRes)&0x80) != 0,
	)

	c.regA = uint8(tmpRes)

	return 1
}

func (c *Cpu) AND() uint8 {
	// Bitwise and
	c.fetch()
	c.regA &= c.fetchedData
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, c.regA&0x80 != 0)
	return 1
}

func (c *Cpu) ASL() uint8 {
	// Arithmetic shift left
	c.fetch()
	tmpRes := uint16(c.fetchedData) << 1
	c.setStatus(FlagC, (tmpRes&0xFF00) > 0)
	c.setStatus(FlagZ, (tmpRes&0x00FF) == 0)
	c.setStatus(FlagN, tmpRes&0x80 != 0)
	if reflect.ValueOf(c.insLookup[c.opcode].addrMode).Pointer() != reflect.ValueOf(c.IMP).Pointer() {
		c.Write(c.addrAbs, uint8(tmpRes))
	} else {
		c.regA = uint8(tmpRes)
	}
	return 0
}

func (c *Cpu) branchLogicHelper() {
	// Many branch use same logic to perform jumping but they had different pre-condition
	c.cyclesLeft++
	c.addrAbs = c.regPC + c.addrRel
	if c.addrAbs&0xFF00 != c.regPC&0xFF00 { // cross page
		c.cyclesLeft++
	}
	c.regPC = c.addrAbs
}

func (c *Cpu) BCC() uint8 {
	// branch if carry clear
	if !c.getStatus(FlagC) {
		c.branchLogicHelper()
	}
	return 0

}

func (c *Cpu) BCS() uint8 {
	// branch if carry set
	if c.getStatus(FlagC) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) BEQ() uint8 {
	// branch if equal
	if c.getStatus(FlagZ) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) BIT() uint8 {
	// logical AND to test the presence of bits in the memory value to set the flags
	// does not keep the result.
	c.fetch()
	tmpRes := c.regA & c.fetchedData
	c.setStatus(FlagZ, (tmpRes&0x00FF) == 0)
	c.setStatus(FlagN, (c.fetchedData&(1<<7)) > 0)
	c.setStatus(FlagV, (c.fetchedData&(1<<6)) > 0)
	return 0
}

func (c *Cpu) BMI() uint8 {
	// branch if negative
	if c.getStatus(FlagN) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) BNE() uint8 {
	// branch if not equal
	if !c.getStatus(FlagZ) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) BPL() uint8 {
	// branch if positive
	if !c.getStatus(FlagN) {
		c.branchLogicHelper()
	}
	return 0

}
func (c *Cpu) BRK() uint8 {
	// Program break
	c.regPC++

	// Set flag and save pc
	c.setStatus(FlagI, true)
	c.Write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC>>8))
	c.regStkPtr--
	c.Write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC))
	c.regStkPtr--

	// save status
	c.setStatus(FlagB, true)
	c.Write(baseStackOffset+uint16(c.regStkPtr), c.regStatus)
	c.regStkPtr--
	c.setStatus(FlagB, false)

	// Set pc, same as interrupt address
	low := c.Read(irqHandlerAddr)
	high := c.Read(irqHandlerAddr + 1)
	c.regPC = (uint16(high) << 8) | uint16(low)

	return 0
}

func (c *Cpu) BVC() uint8 {
	// branch if overflow clear
	if !c.getStatus(FlagV) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) BVS() uint8 {
	// branch if overflow set
	if c.getStatus(FlagV) {
		c.branchLogicHelper()
	}
	return 0
}

func (c *Cpu) CLC() uint8 {
	// Clear carry
	c.setStatus(FlagC, false)
	return 0
}

func (c *Cpu) CLD() uint8 {
	// Clear decimal
	c.setStatus(FlagD, false)
	return 0
}

func (c *Cpu) CLI() uint8 {
	// Clear disable interrupt
	c.setStatus(FlagI, false)
	return 0
}

func (c *Cpu) CLV() uint8 {
	// Clear overflow
	c.setStatus(FlagV, false)
	return 0
}

func (c *Cpu) CMP() uint8 {
	// compare accum
	c.fetch()
	tmpRes := uint16(c.regA) - uint16(c.fetchedData)

	// Check flag
	c.setStatus(FlagC, c.regA >= c.fetchedData)
	c.setStatus(FlagZ, (tmpRes&0xFF) == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)

	return 1
}

func (c *Cpu) CPX() uint8 {
	// compare X
	c.fetch()
	tmpRes := uint16(c.regX) - uint16(c.fetchedData)

	// Check flag
	c.setStatus(FlagC, c.regX >= c.fetchedData)
	c.setStatus(FlagZ, (tmpRes&0xFF) == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)

	return 1
}
func (c *Cpu) CPY() uint8 {
	// compare Y
	c.fetch()
	tmpRes := uint16(c.regY) - uint16(c.fetchedData)

	// Check flag
	c.setStatus(FlagC, c.regY >= c.fetchedData)
	c.setStatus(FlagZ, (tmpRes&0xFF) == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)

	return 1

}
func (c *Cpu) DEC() uint8 {
	// Decrement value at memory location
	c.fetch()
	tmpRes := c.fetchedData - 1
	c.Write(c.addrAbs, tmpRes)
	c.setStatus(FlagZ, tmpRes == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)
	return 0
}
func (c *Cpu) DEX() uint8 {
	// Decrement x
	c.regX--
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 0
}
func (c *Cpu) DEY() uint8 {
	// Decrement y
	c.regY--
	c.setStatus(FlagZ, c.regY == 0)
	c.setStatus(FlagN, (c.regY&0x80) != 0)
	return 0
}

func (c *Cpu) EOR() uint8 {
	// Bitwise XOR
	c.fetch()
	c.regA = c.regA ^ c.fetchedData
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, c.regA&0x80 != 0)
	return 1
}

func (c *Cpu) INC() uint8 {
	// Increment value at memory location
	c.fetch()
	tmpRes := c.fetchedData + 1
	c.Write(c.addrAbs, tmpRes)
	c.setStatus(FlagZ, tmpRes == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)
	return 0
}

func (c *Cpu) INX() uint8 {
	// Increment x
	c.regX++
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 0
}

func (c *Cpu) INY() uint8 {
	// Increment y
	c.regY++
	c.setStatus(FlagZ, c.regY == 0)
	c.setStatus(FlagN, (c.regY&0x80) != 0)
	return 0
}

func (c *Cpu) JMP() uint8 {
	// Direct jump
	c.regPC = c.addrAbs
	return 0
}

func (c *Cpu) JSR() uint8 {
	// Jump to subroutine, save pc
	c.regPC--
	c.Write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC>>8))
	c.regStkPtr--
	c.Write(baseStackOffset+uint16(c.regStkPtr), uint8(c.regPC))
	c.regStkPtr--

	c.regPC = c.addrAbs
	return 0
}

func (c *Cpu) LDA() uint8 {
	// LoadToRam to accum
	c.fetch()
	c.regA = c.fetchedData
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, (c.regA&0x80) != 0)
	return 1
}

func (c *Cpu) LDX() uint8 {
	// LoadToRam to X
	c.fetch()
	c.regX = c.fetchedData
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 1
}

func (c *Cpu) LDY() uint8 {
	// LoadToRam to Y
	c.fetch()
	c.regY = c.fetchedData
	c.setStatus(FlagZ, c.regY == 0)
	c.setStatus(FlagN, (c.regY&0x80) != 0)
	return 1
}

func (c *Cpu) LSR() uint8 {
	// Logical Shift Right
	c.fetch()
	c.setStatus(FlagC, (c.fetchedData&1) > 0)
	tmpRes := c.fetchedData >> 1
	c.setStatus(FlagZ, tmpRes == 0)
	c.setStatus(FlagN, tmpRes&0x80 != 0)
	if reflect.ValueOf(c.insLookup[c.opcode].addrMode).Pointer() != reflect.ValueOf(c.IMP).Pointer() {
		c.Write(c.addrAbs, tmpRes)
	} else {
		c.regA = tmpRes
	}
	return 0
}

func (c *Cpu) NOP() uint8 {
	// not all NOPs are equal: https://wiki.nesdev.com/w/index.php/CPU_unofficial_opcodes
	// Some NOP will use non-immediate mode address and consume extra byte after nop
	// TODO: will add more based on game compatibility, and ultimately, I'd like to cover all illegal opcodes too
	switch c.opcode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return 1
	}
	return 0
}

func (c *Cpu) ORA() uint8 {
	// Bitwise OR
	c.fetch()
	c.regA |= c.fetchedData
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, (c.regA&0x80) != 0)
	return 1
}

func (c *Cpu) PHA() uint8 {
	// Push accumulator to stack
	c.Write(baseStackOffset+uint16(c.regStkPtr), c.regA)
	c.regStkPtr--
	return 0
}

func (c *Cpu) PHP() uint8 {
	// Push status to stack
	c.Write(baseStackOffset+uint16(c.regStkPtr), c.regStatus|FlagB|FlagU)
	c.regStkPtr--
	return 0
}

func (c *Cpu) PLA() uint8 {
	// Pop stack to accumulator, set flag
	c.regStkPtr++
	c.regA = c.Read(baseStackOffset + uint16(c.regStkPtr))
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, (c.regA&0x80) > 0)
	return 0
}

func (c *Cpu) PLP() uint8 {
	// Pop stack to status, set flag
	c.regStkPtr++
	c.regStatus = c.Read(baseStackOffset + uint16(c.regStkPtr))
	c.setStatus(FlagB, false)
	c.setStatus(FlagU, true)
	return 0
}

func (c *Cpu) ROL() uint8 {
	// Rotate left, similar to shift
	c.fetch()
	tmpRes := uint16(c.fetchedData) << 1
	if c.getStatus(FlagC) {
		tmpRes |= 1
	}
	c.setStatus(FlagC, (tmpRes&0xFF00) > 0)
	c.setStatus(FlagZ, (tmpRes&0x00FF) == 0)
	c.setStatus(FlagN, tmpRes&0x80 != 0)
	if reflect.ValueOf(c.insLookup[c.opcode].addrMode).Pointer() != reflect.ValueOf(c.IMP).Pointer() {
		c.Write(c.addrAbs, uint8(tmpRes))
	} else {
		c.regA = uint8(tmpRes)
	}
	return 0
}

func (c *Cpu) ROR() uint8 {
	// Rotate right, similar to shift
	c.fetch()
	tmpRes := uint16(c.fetchedData) >> 1
	if c.getStatus(FlagC) {
		tmpRes |= 1 << 7
	}
	c.setStatus(FlagC, (c.fetchedData&1) > 0)
	c.setStatus(FlagZ, (tmpRes&0x00FF) == 0)
	c.setStatus(FlagN, tmpRes&0x80 != 0)
	if reflect.ValueOf(c.insLookup[c.opcode].addrMode).Pointer() != reflect.ValueOf(c.IMP).Pointer() {
		c.Write(c.addrAbs, uint8(tmpRes))
	} else {
		c.regA = uint8(tmpRes)
	}
	return 0
}

func (c *Cpu) RTI() uint8 {
	// Return from interrupt
	// Restore state, reverse logic from handling interrupt
	c.regStkPtr++
	c.regStatus = c.Read(baseStackOffset + uint16(c.regStkPtr))
	c.setStatus(FlagU, true)

	c.regStkPtr++
	c.regPC = uint16(c.Read(baseStackOffset + uint16(c.regStkPtr)))
	c.regStkPtr++
	c.regPC |= uint16(c.Read(baseStackOffset+uint16(c.regStkPtr))) << 8

	return 0
}

func (c *Cpu) RTS() uint8 {
	// Return from subroutine
	// Restore state
	c.regStkPtr++
	c.regPC = uint16(c.Read(baseStackOffset + uint16(c.regStkPtr)))
	c.regStkPtr++
	c.regPC |= uint16(c.Read(baseStackOffset+uint16(c.regStkPtr))) << 8
	c.regPC++
	return 0
}

func (c *Cpu) SBC() uint8 {
	// Subtraction, almost identical to addition
	c.fetch()
	invertVal := uint16(c.fetchedData) ^ 0xFF // the only different line!
	tmpRes := uint16(c.regA) + invertVal
	if c.getStatus(FlagC) {
		tmpRes++
	}

	// Check flag
	c.setStatus(FlagC, tmpRes > 0xFF)
	c.setStatus(FlagZ, (tmpRes&0xFF) == 0)
	c.setStatus(FlagN, (tmpRes&0x80) != 0)
	c.setStatus(FlagV,
		((tmpRes^invertVal)&(uint16(c.regA)^tmpRes)&0x80) != 0,
	)

	c.regA = uint8(tmpRes)

	return 1
}

func (c *Cpu) SEC() uint8 {
	// Set carry
	c.setStatus(FlagC, true)
	return 0
}

func (c *Cpu) SED() uint8 {
	// Set decimal
	c.setStatus(FlagD, true)
	return 0
}

func (c *Cpu) SEI() uint8 {
	// Set disable interrupt
	c.setStatus(FlagI, true)
	return 0
}

func (c *Cpu) STA() uint8 {
	// Store accum at address
	c.Write(c.addrAbs, c.regA)
	return 0

}
func (c *Cpu) STX() uint8 {
	// Store x at address
	c.Write(c.addrAbs, c.regX)
	return 0

}

func (c *Cpu) STY() uint8 {
	// Store y at address
	c.Write(c.addrAbs, c.regY)
	return 0
}

func (c *Cpu) TAX() uint8 {
	// Transfer a -> x
	c.regX = c.regA
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 0
}

func (c *Cpu) TAY() uint8 {
	// Transfer a -> y
	c.regY = c.regA
	c.setStatus(FlagZ, c.regY == 0)
	c.setStatus(FlagN, (c.regY&0x80) != 0)
	return 0
}

func (c *Cpu) TSX() uint8 {
	// Transfer stkPtr -> x
	c.regX = c.regStkPtr
	c.setStatus(FlagZ, c.regX == 0)
	c.setStatus(FlagN, (c.regX&0x80) != 0)
	return 0
}

func (c *Cpu) TXA() uint8 {
	// Transfer x -> a
	c.regA = c.regX
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, (c.regA&0x80) != 0)
	return 0
}

func (c *Cpu) TXS() uint8 {
	// Transfer x -> stkPtr, note no status testing
	c.regStkPtr = c.regX
	return 0
}

func (c *Cpu) TYA() uint8 {
	// Transfer y -> a
	c.regA = c.regY
	c.setStatus(FlagZ, c.regA == 0)
	c.setStatus(FlagN, (c.regA&0x80) != 0)
	return 0
}

// XXX for illegal opcode! Defined by ourselves
func (c *Cpu) XXX() uint8 {
	mlog.L.Fatal("Encountered invalid opcode")
	return 0
}
