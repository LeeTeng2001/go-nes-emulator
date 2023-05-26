package bus

import (
	"nes_emulator/disk"
	"nes_emulator/mlog"
)

const (
	RamPhysicalSize      = 0x0800
	RamAccessMask        = 0x07FF
	RamAccessSize        = 0x2000
	PpuAccessRegionEnd   = 0x4000
	PpuAccessMask        = 0x07
	InputWriteBoundStart = 0x4016
	InputWriteBoundEnd   = 0x4017
)

// TODO: Barrier to prevent unauthorised read/write by other device?
// TODO: Check disk is valid when accessing

// Bus connect all components, device is defined as interface
// to prevent circular dependencies and support custom device implementation
type Bus struct {
	cpu      CpuDevice
	ppu      PpuDevice
	disk     *disk.NesDisk
	ram      [RamPhysicalSize]uint8
	sysClock uint64
	// Input bits for two controller at current frame and it's stored state
	controllerInput [2]uint8
	controllerState [2]uint8
}

func New(Cpu CpuDevice, Ppu PpuDevice) *Bus {
	b := Bus{
		cpu:      Cpu,
		ppu:      Ppu,
		sysClock: 0,
	}
	mlog.L.Infof("Bus is initialised")
	return &b
}

func (b *Bus) InsertDisk(nesDisk *disk.NesDisk) {
	b.disk = nesDisk
	if b.ppu != nil {
		b.ppu.ConnectDisk(nesDisk)
	} else {
		mlog.L.Warnf("Ppu is empty when inserting disk")
	}
}

func (b *Bus) Reset() {
	mlog.L.Info("Resetting bus")
	b.cpu.Reset()
	b.ppu.Reset()
	b.sysClock = 0
}

func (b *Bus) UpdateControllerInputBits(controllerIdx uint8, val uint8) {
	if controllerIdx > 1 {
		mlog.L.Fatal("Only 2 controller is supported")
	}
	b.controllerInput[controllerIdx] = val
}

func (b *Bus) Clock() {
	// The running frequency is controlled by whatever calls this function.

	// The fastest clock frequency the digital system cares is ppu
	b.ppu.Clock()

	// The CPU runs 3 times slower than the PPU
	if b.sysClock%3 == 0 {
		b.cpu.Clock()
	}

	// Check nmi in ppu and initiate it in cpu
	if b.ppu.CheckNmiAndTurnOff() {
		b.cpu.Nmi()
	}

	// Update sys clock
	b.sysClock++
}

// LoadToRam : directly load a range of data into ram, only for debugging
func (b *Bus) LoadToRam(data []uint8, startAddr int) {
	for offset, dataByte := range data {
		addr := startAddr + offset
		if addr >= 0 && addr <= 0xFFFF {
			b.ram[addr] = dataByte
		} else {
			mlog.L.Fatal("LoadToRam memory out of range!")
		}
	}
}
