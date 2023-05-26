package nes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"io"
	"io/fs"
	"nes_emulator/bus"
	"nes_emulator/cpu6502"
	"nes_emulator/disk"
	"nes_emulator/mlog"
	"nes_emulator/olc2c02"
)

// Game implement the gui and is the glue between user and devices
type Game struct {
	// Nes components, put ppu here to directly access the display buffer
	bus     *bus.Bus
	ppu     *olc2c02.Ppu
	nesDisk *disk.NesDisk
	// Controller
	controllerMap map[ebiten.Key]uint8
	// Screen info
	screenW int
	screenH int
	// Draw medium
	screenImg *ebiten.Image
	screenPT0 *ebiten.Image
	screenPT1 *ebiten.Image
	// Draw options
	diOptMain          *ebiten.DrawImageOptions
	diOptPatternTable0 *ebiten.DrawImageOptions
	diOptPatternTable1 *ebiten.DrawImageOptions
	// State
	selectedPaletteId uint8
	textFont          font.Face
}

func New() *Game {
	// Assemble different part of hardware
	//cpu := cpu6502.NewDebug() // TODO: Remove debug
	cpu := cpu6502.New()
	ppu := olc2c02.New()
	nesBus := bus.New(cpu, ppu)
	cpu.ConnectBus(nesBus)
	ppu.ConnectBus(nesBus)

	// Setup draw buffer, screen
	g := &Game{
		screenW:            256 + 20 + 256/2,
		screenH:            240 + 16,
		bus:                nesBus,
		ppu:                ppu,
		controllerMap:      make(map[ebiten.Key]uint8),
		screenImg:          ebiten.NewImage(256, 240),
		screenPT0:          ebiten.NewImage(128, 128),
		screenPT1:          ebiten.NewImage(128, 128),
		diOptMain:          &ebiten.DrawImageOptions{},
		diOptPatternTable0: &ebiten.DrawImageOptions{},
		diOptPatternTable1: &ebiten.DrawImageOptions{},
	}
	// Setup controller map
	g.initControllerMap()
	// Setup draw region
	g.diOptMain.GeoM.Translate(0, 16/2)
	g.diOptPatternTable0.GeoM.Translate(256+20, 0)
	g.diOptPatternTable1.GeoM.Translate(256+20, (240)/2+16)

	// Initialise draw font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		mlog.L.Fatal(err)
	}
	g.textFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     100,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		mlog.L.Fatal(err)
	}

	// TODO: Setup instruction
	mlog.L.Info("Draw screen is initialised")
	return g
}

func (g *Game) Run() {
	mlog.L.Info("Running game loop...")
	ebiten.SetWindowSize(g.screenW*2, g.screenH*2)
	ebiten.SetWindowTitle("NES Emulator")
	if err := ebiten.RunGame(g); err != nil {
		mlog.L.Fatal(err)
	}
}

func (g *Game) Update() error {
	// Check drop file input to refresh nes
	if files := ebiten.DroppedFiles(); files != nil {
		if err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				mlog.L.Fatalf("Not possible errors encounter during walkdir: %s", err)
			}
			if d.IsDir() { // skip dir
				return nil
			}
			mlog.L.Printf("Name: %s, Path: %s", d.Name(), path)
			data, err := fs.ReadFile(files, path) // read file content
			if err != nil {
				mlog.L.Fatalf("Error occurs when reading content of: %s", path)
			}
			f, err := files.Open(path) // open file
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()
			// check valid nes file
			newNesDisk := disk.NewFromBytes(data)
			if newNesDisk != nil {
				newNesDisk.PrintInfo()
				g.bus.InsertDisk(newNesDisk)
				g.bus.Reset()
				g.nesDisk = newNesDisk
				return io.EOF
			}
			return nil
		}); err != nil {

		}
	}
	// Valid game is running
	if g.nesDisk != nil {
		// Update controller
		g.updateControllerState()
		// Run simulation for whole frame
		// TODO: Select palette, fix clock
		for !g.ppu.FrameCompleteAndTurnOff() {
			g.bus.Clock()
		}
		// Check for debug input
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.selectedPaletteId = (g.selectedPaletteId + 1) % 8
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			mlog.L.Infof("TPS: %0.3f", ebiten.ActualTPS())
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.nesDisk == nil {
		g.DrawMainMenu(screen)
	} else {
		g.DrawGame(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenW, g.screenH
}
