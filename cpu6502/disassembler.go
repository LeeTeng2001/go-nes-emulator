package cpu6502

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Turn a range of memory into human readable format
func (c *Cpu) disassembleMemRange(startAddr, endAddr uint16, skipNop bool) string {
	var output strings.Builder

	lineOutput := ""
	for startAddr <= endAddr {
		lineOutput, startAddr = c.disassembleMemAtAddr(startAddr)
		output.WriteString(lineOutput)
		output.WriteByte('\n')
	}

	fmt.Println(output.String())

	return output.String()
}

func (c *Cpu) disassembleMemAtAddr(addr uint16) (ans string, nextAddr uint16) {
	// Example line: C000  4C F5 C5  JMP $C5F5
	var output strings.Builder

	// read instruction and base operand
	output.WriteString(fmt.Sprintf("%04X  ", addr))
	opcode := c.read(addr)
	addr++
	inst := c.insLookup[opcode]
	output.WriteString(fmt.Sprintf("%02X ", opcode))

	addrModeStr := runtime.FuncForPC(reflect.ValueOf(inst.addrMode).Pointer()).Name()
	addrModeStr = addrModeStr[len(addrModeStr)-6 : len(addrModeStr)-3]

	// form the instruction based upon its addressing mode
	infoString := ""
	switch reflect.ValueOf(inst.addrMode).Pointer() {
	case reflect.ValueOf(c.IMM).Pointer(), reflect.ValueOf(c.ZP0).Pointer(),
		reflect.ValueOf(c.ZPX).Pointer(), reflect.ValueOf(c.ZPY).Pointer(),
		reflect.ValueOf(c.IZX).Pointer(), reflect.ValueOf(c.IZY).Pointer():
		lowData := c.read(addr)
		addr++
		infoString = fmt.Sprintf("$%02X", lowData)
		output.WriteString(fmt.Sprintf("%02X    ", lowData))
	case reflect.ValueOf(c.ABS).Pointer(), reflect.ValueOf(c.ABX).Pointer(),
		reflect.ValueOf(c.ABY).Pointer(), reflect.ValueOf(c.IND).Pointer():
		lowData := c.read(addr)
		addr++
		highData := c.read(addr)
		addr++
		infoString = fmt.Sprintf("$%02X%02X", highData, lowData)
		output.WriteString(fmt.Sprintf("%02X %02X ", lowData, highData))
	case reflect.ValueOf(c.REL).Pointer(): // handle redirection maybe
		lowData := c.read(addr)
		addr++
		output.WriteString(fmt.Sprintf("%02X    ", lowData))
	case reflect.ValueOf(c.NOP).Pointer(): // NOP
		lowData := c.read(addr)
		addr++
		output.WriteString(fmt.Sprintf("%02X    ", lowData))
	default:
		output.WriteString("      ")
	}

	// Final output
	output.WriteString(fmt.Sprintf(" %s (%s) %-8s", inst.name, addrModeStr, infoString))
	return output.String(), addr
}

func (c *Cpu) getRegStateStr() string {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		c.regA, c.regX, c.regY, c.regStatus, c.regStkPtr)
}
