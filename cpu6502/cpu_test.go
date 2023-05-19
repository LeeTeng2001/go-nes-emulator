package cpu6502

import (
	"bufio"
	"github.com/charmbracelet/log"
	"nes_emulator/bus"
	"nes_emulator/loader"
	"nes_emulator/mlog"
	"os"
	"testing"
)

func compareLog() bool {
	// Open files
	nesTestLog, err := os.Open("tests/nesttest.log")
	if err != nil {
		log.Fatal("Cannot open tests/nesttest.log")
	}
	cpuLog, err := os.Open(cpuLogFile)
	if err != nil {
		log.Fatal("Cannot open " + cpuLogFile)
	}

	// Scanner
	nesTestScanner := bufio.NewScanner(nesTestLog)
	nesTestScanner.Split(bufio.ScanLines)
	cpuScanner := bufio.NewScanner(cpuLog)
	cpuScanner.Split(bufio.ScanLines)

	curLine := 1
	for cpuScanner.Scan() && nesTestScanner.Scan() {
		cpuLine := cpuScanner.Text()
		nesTestLine := nesTestScanner.Text()
		// TODO: Hardcode position and other information
		cpuFrontInfo := cpuLine[23:38]
		cpuBackInfo := cpuLine[58:]
		nesTestFrontInfo := nesTestLine[:15]
		nesTestBackInfo := nesTestLine[48:73]

		if cpuFrontInfo != nesTestFrontInfo {
			log.Errorf("Front info mismatch at line %v", curLine)
			log.Errorf("cpu    : %s", cpuFrontInfo)
			log.Errorf("nestest: %s", nesTestFrontInfo)
			return false
		}

		if cpuBackInfo != nesTestBackInfo {
			log.Errorf("Back info mismatch at line %v", curLine)
			log.Errorf("cpu    : %s", cpuBackInfo)
			log.Errorf("nestest: %s", nesTestBackInfo)
			return false
		}

		curLine++
	}

	return true
}

func TestNesTest(t *testing.T) {
	newBus := bus.New()
	cpu := NewDebug()
	cpu.ConnectBus(newBus)

	// Load nes, write twice since we don't have mapper
	// More info: https://github.com/PyAndy/Py3NES/issues/1
	nesFile := loader.New("tests/nestest.nes")
	newBus.LoadToRam(nesFile.PrgRomData, 0x8000)
	newBus.LoadToRam(nesFile.PrgRomData, 0xC000)

	// Set execution start point for nestest
	cpu.reset()
	cpu.regPC = 0xC000
	//cpu.disassembleMemRange(0xC000, 0xC005, false)

	// Run, after this value is unofficial opcode which is discouraged
	// https://www.nesdev.org/wiki/Programming_with_unofficial_opcodes
	for i := 0; i < 15850; i++ {
		cpu.Clock()
	}

	// check for error as defined by nestest
	// Doc: https://www.qmtpro.com/~nes/misc/nestest.txt
	errCodeLow := cpu.Read(0x02)
	errCodeHigh := cpu.Read(0x03)
	mlog.L.Infof("Err code low: %02X", errCodeLow)
	mlog.L.Infof("Err code high: %02X", errCodeHigh)

	if errCodeLow != 0 || errCodeHigh != 0 {
		t.Errorf("One of the error code is not zero!")
	}

	// Compare log
	if !compareLog() {
		t.Errorf("Log output mismatch!!")
	}
}
