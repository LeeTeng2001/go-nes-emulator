package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"nes_emulator/bus"
	"nes_emulator/cpu6502"
	"nes_emulator/mlog"
	"nes_emulator/olc2c02"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
	screen.Set(50, 50, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	// Assemble different part of hardware
	cpu := cpu6502.New()
	ppu := olc2c02.New()
	newBus := bus.New(cpu, ppu)
	cpu.ConnectBus(newBus)
	ppu.ConnectBus(newBus)

	//for {
	//	cpu.Clock()
	//}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		mlog.L.Fatal(err)
	}

}
