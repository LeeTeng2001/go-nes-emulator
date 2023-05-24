package disk

import (
	"nes_emulator/disk/mapper"
	"nes_emulator/mlog"
)

// INES: https://www.nesdev.org/wiki/INES

const (
	HeaderSize     = 16
	PrgRomSizeUnit = 0x4000
	ChrRomSizeUnit = 0x2000
)

// NesDisk contains the loaded nes file content
// and handle interaction to the main / cpu bus
type NesDisk struct {
	FormatIsNes2 bool
	// Header (16 bytes, last 6 bytes unused)
	PrgTotalBank uint8
	ChrTotalBank uint8
	mapper1      uint8
	mapper2      uint8
	flag8        uint8
	flag9        uint8
	flag10       uint8
	// Data
	PrgRomData []byte
	ChrRomData []byte
	// Mapper info and corresponding handler
	nMapperId     uint8
	mapperHandler mapper.IMapper
}

func (n *NesDisk) PrintInfo() {
	format := "iNES"
	if n.FormatIsNes2 {
		format = "NES2.0"
	}

	mlog.L.Infof("Nes format: %s", format)
	mlog.L.Infof("Header:")
	mlog.L.Infof("\tPRG bank amt (each 16KB): %v", n.PrgTotalBank)
	mlog.L.Infof("\tCHR bank amt (each 8KB): %v", n.ChrTotalBank)
	mlog.L.Infof("\tMapper ID: %v", n.nMapperId)
}
