package nes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	text.Draw(screen, "Drag and drop a nesfile to run", g.textFont, 80, 80, color.White)
	text.Draw(screen, "P: Switch selected palette", g.textFont, 80, 120, color.White)
	text.Draw(screen, "F: Print true FPS in console", g.textFont, 80, 140, color.White)
}

func (g *Game) DrawGame(screen *ebiten.Image) {
	// Background fill
	screen.Fill(color.RGBA{
		R: 76,
		G: 52,
		B: 235,
		A: 255,
	})

	// TODO: Support upscaling
	// Draw main screen and 2 pattern tables
	g.screenImg.WritePixels(g.ppu.GetScreenPixels())
	g.screenPT0.WritePixels(g.ppu.GetPatternTable(g.selectedPaletteId, 0))
	g.screenPT1.WritePixels(g.ppu.GetPatternTable(g.selectedPaletteId, 1))
	screen.DrawImage(g.screenImg, g.diOptMain)
	screen.DrawImage(g.screenPT0, g.diOptPatternTable0)
	screen.DrawImage(g.screenPT1, g.diOptPatternTable1)

	// Draw palettes
	allColors := g.ppu.GetAllColorPalettes()
	startX := float32(256 + 5)
	startY := float32(10)
	areaHeight := float32(256-20-5*7) / 8
	blockHeight := areaHeight / 4
	areaIdx := -1
	for i, singleColor := range allColors {
		if i%4 == 0 {
			areaIdx++
			// Draw selected background
			if int(g.selectedPaletteId) == areaIdx {
				vector.DrawFilledRect(screen, startX-1, startY+float32(areaIdx)*(areaHeight+5)-1,
					10+2, areaHeight+2, color.White, false)
			}
		}
		localBlockOffset := float32(i%4) * blockHeight
		vector.DrawFilledRect(screen, startX, startY+float32(areaIdx)*(areaHeight+5)+localBlockOffset,
			10, blockHeight, singleColor, false)
	}
}
