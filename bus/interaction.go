package bus

func (b *Bus) CWrite(addr uint16, data uint8) {
	if b.disk.CWrite(addr, data) {
		// The cartridge "sees all" and has the facility to veto
		// the propagation of the bus transaction if it requires.
		// This allows the cartridge to map any address to some
		// other data, including the facility to divert transactions
		// with other physical devices. The NES does not do this!!!!!!
		// This allow "custom" hardware behaviour to the NES in the future!
	} else if addr < RamAccessSize {
		// Ram address range, it mirrors its physical region 4 times
		b.ram[addr&RamAccessMask] = data
	} else if addr >= RamAccessSize && addr < PpuAccessRegionEnd {
		// Ppu address range, it only has 8 registers and it's being mirrored
		// https://www.nesdev.org/wiki/PPU_registers
		b.ppu.CWrite(addr&PpuAccessMask, data)
	} else if (addr >= 0x4000 && addr <= 0x4013) || addr == 0x4015 || addr == 0x4017 { // apu
		b.apu.CWrite(addr, data)
	} else if addr == 0x4014 { // DMA started, transfer one full page
		b.dmaPage = data
		b.dmaAddr = 0
		b.dmaTransferIsOn = true
		b.dmaDummyIsOn = true
	} else if addr >= InputWriteBoundStart && addr <= InputWriteBoundEnd {
		// Poll controller to takes a snapshot of current state: https://www.nesdev.org/wiki/Controller_reading
		// 2 bits for two controller
		b.controllerState[addr&1] = b.controllerInput[addr&1]
	}
}

func (b *Bus) CRead(addr uint16) (data uint8) {
	if b.disk.CRead(addr, &data) {
		// Same purpose for custom modification
		return data
	} else if addr < RamAccessSize {
		return b.ram[addr&RamAccessMask]
	} else if addr >= RamAccessSize && addr < PpuAccessRegionEnd {
		return b.ppu.CRead(addr & PpuAccessMask)
	} else if addr >= InputWriteBoundStart && addr <= InputWriteBoundEnd {
		// Read the high bit of stored selected controller state and shift it for further read
		highBit := b.controllerState[addr&1] & 0x80
		b.controllerState[addr&1] <<= 1
		if highBit != 0 {
			return 1
		} else {
			return 0
		}
	}
	return 0
}
