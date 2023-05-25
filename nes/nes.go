package nes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"nes_emulator/bus"
	"nes_emulator/cpu6502"
	"nes_emulator/disk"
	"nes_emulator/mlog"
	"nes_emulator/olc2c02"
)

// Game implement the gui and is the glue between user and devices
type Game struct {
	// Nes components
	bus *bus.Bus
	ppu *olc2c02.Ppu
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
	// Input
	keys []ebiten.Key
}

func New(nesDiskPath string) *Game {
	// Assemble different part of hardware
	//cpu := cpu6502.NewDebug() // TODO: Remove debug
	cpu := cpu6502.New()
	ppu := olc2c02.New()
	nesBus := bus.New(cpu, ppu)
	cpu.ConnectBus(nesBus)
	ppu.ConnectBus(nesBus)
	nesFile := disk.New(nesDiskPath)
	//nesFile.PrintInfo()
	nesBus.InsertDisk(nesFile)
	nesBus.Reset()

	// Setup draw buffer, screen
	g := &Game{
		screenW:            256 + 20 + 256/2,
		screenH:            240 + 16,
		bus:                nesBus,
		ppu:                ppu,
		screenImg:          ebiten.NewImage(256, 240),
		screenPT0:          ebiten.NewImage(128, 128),
		screenPT1:          ebiten.NewImage(128, 128),
		diOptMain:          &ebiten.DrawImageOptions{},
		diOptPatternTable0: &ebiten.DrawImageOptions{},
		diOptPatternTable1: &ebiten.DrawImageOptions{},
	}
	// Setup draw region
	g.diOptMain.GeoM.Translate(0, 16/2)
	g.diOptPatternTable0.GeoM.Translate(256+20, 0)
	g.diOptPatternTable1.GeoM.Translate(256+20, (240)/2+16)

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
	// TODO: Select palette, fix clock
	for !g.ppu.FrameCompleteAndTurnOff() {
		g.bus.Clock()
	}
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Fill(color.RGBA{
		R: 15,
		G: 25,
		B: 100,
		A: 255,
	})
	g.screenImg.WritePixels(g.ppu.GetScreenPixels())
	g.screenPT0.WritePixels(g.ppu.GetPatternTable(0))
	g.screenPT1.WritePixels(g.ppu.GetPatternTable(1))
	screen.DrawImage(g.screenImg, g.diOptMain)
	screen.DrawImage(g.screenPT0, g.diOptPatternTable0)
	screen.DrawImage(g.screenPT1, g.diOptPatternTable1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenW, g.screenH
}
