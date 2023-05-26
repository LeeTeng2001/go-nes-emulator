package disk

import (
	"nes_emulator/disk/mapper"
	"nes_emulator/mlog"
	"os"
)

func New(filepath string) *NesDisk {
	data, err := os.ReadFile(filepath)
	if err != nil {
		mlog.L.Fatal("Failed to read: " + filepath)
	}
	return NewFromBytes(data)
}

func NewFromBytes(data []byte) *NesDisk {
	nesFile := NesDisk{}

	// Check valid
	if len(data) < 4 {
		mlog.L.Error("Data is too small for a nes file")
		return nil
	}
	if data[0] == 'N' && data[1] == 'E' && data[2] == 'S' && data[3] == 0x1A {
		nesFile.FormatIsNes2 = false
		// Check nes2 format
		if data[7]&0x0C == 0x08 {
			nesFile.FormatIsNes2 = true
		}
	} else {
		mlog.L.Error("Data is not a nes file")
		return nil
	}

	// Load, TODO assume it's correct
	nesFile.PrgTotalBank = data[4]
	nesFile.ChrTotalBank = data[5]
	nesFile.mapper1 = data[6]
	nesFile.mapper2 = data[7]
	nesFile.flag8 = data[8]
	nesFile.flag9 = data[9]
	nesFile.flag10 = data[10]

	// Skip trainer data
	trainerSize := 0
	if nesFile.mapper1&0x04 != 0 {
		trainerSize = 512
	}

	// Determine mapper id and mirror info
	nesFile.nMapperId = (nesFile.mapper2 & 0xF0) | (nesFile.mapper1 >> 4)
	if nesFile.mapper1&1 == 0 {
		nesFile.MirrorHorizontal = true
	}

	// Load ROM data
	prgActualSize := PrgRomSizeUnit * int(nesFile.PrgTotalBank)
	chrActualSize := ChrRomSizeUnit * int(nesFile.ChrTotalBank)
	nesFile.PrgRomData = data[HeaderSize+trainerSize : HeaderSize+trainerSize+prgActualSize]
	nesFile.ChrRomData = data[HeaderSize+trainerSize+prgActualSize : HeaderSize+trainerSize+prgActualSize+chrActualSize]

	// Load correspond mapper handler
	switch nesFile.nMapperId {
	case 0:
		nesFile.mapperHandler = mapper.NewM0(nesFile.PrgTotalBank, nesFile.ChrTotalBank)
	default:
		mlog.L.Fatalf("Encounter mapper id %d without corresponding mapper!\n", nesFile.nMapperId)
	}

	return &nesFile
}

// TODO: Load to memory directly
