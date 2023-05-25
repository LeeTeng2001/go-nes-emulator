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
	}
	return 0
}
