package loader

import (
	"nes_emulator/mlog"
	"os"
)

// INES: https://www.nesdev.org/wiki/INES

const (
	HeaderSize     = 16
	PrgRomSizeUnit = 0x4000
	ChrRomSizeUnit = 0x2000
)

type NesFile struct {
	FormatIsNes2 bool
	// Header (16 bytes, last 6 bytes unused)
	PrgRomSize uint8
	ChrRomSize uint8
	flag6      uint8
	flag7      uint8
	flag8      uint8
	flag9      uint8
	flag10     uint8
	// Data
	PrgRomData []byte
	ChrRomData []byte
}

func New(filepath string) *NesFile {
	data, err := os.ReadFile(filepath)
	if err != nil {
		mlog.L.Fatal("Failed to read: " + filepath)
	}
	nesFile := NesFile{}

	// Check valid
	if len(data) < 4 {
		mlog.L.Fatal("file is not a nes file: " + filepath)
	}
	if data[0] == 'N' && data[1] == 'E' && data[2] == 'S' && data[3] == 0x1A {
		nesFile.FormatIsNes2 = false
		// Check nes2
		if data[7]&0x0C == 0x08 {
			nesFile.FormatIsNes2 = true
		}
	} else {
		mlog.L.Fatal("file is not a nes file: " + filepath)
	}

	// Load, TODO assume it's correct
	nesFile.PrgRomSize = data[4]
	nesFile.ChrRomSize = data[5]
	nesFile.flag6 = data[6]
	nesFile.flag7 = data[7]
	nesFile.flag8 = data[8]
	nesFile.flag9 = data[9]
	nesFile.flag10 = data[10]

	// TODO: Optional trainer data

	// ROM data
	prgActualSize := PrgRomSizeUnit * int(nesFile.PrgRomSize)
	chrActualSize := ChrRomSizeUnit * int(nesFile.ChrRomSize)
	nesFile.PrgRomData = data[HeaderSize : HeaderSize+prgActualSize]
	nesFile.ChrRomData = data[HeaderSize+prgActualSize : HeaderSize+prgActualSize+chrActualSize]

	return &nesFile
}

func (n *NesFile) PrintInfo() {
	format := "iNES"
	if n.FormatIsNes2 {
		format = "NES2.0"
	}

	mlog.L.Infof("Nes format: %s", format)
	mlog.L.Infof("Header:")
	mlog.L.Infof("\tPRG Rom size (in 16KB): %v", n.PrgRomSize)
	mlog.L.Infof("\tCHR Rom size (in 8KB): %v", n.ChrRomSize)
}
