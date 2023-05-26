package nes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) initControllerMap() {
	// Map input to corresponding bit
	// TODO: Support 2nd controllers, and possible other controller input type, and key customisation
	g.controllerMap[ebiten.KeyX] = 0x80
	g.controllerMap[ebiten.KeyZ] = 0x40
	g.controllerMap[ebiten.KeyA] = 0x20
	g.controllerMap[ebiten.KeyS] = 0x10
	g.controllerMap[ebiten.KeyUp] = 0x08
	g.controllerMap[ebiten.KeyDown] = 0x04
	g.controllerMap[ebiten.KeyLeft] = 0x02
	g.controllerMap[ebiten.KeyRight] = 0x01
}

func (g *Game) updateControllerState() {
	inputBit := uint8(0)
	for k, bitCode := range g.controllerMap {
		if inpututil.KeyPressDuration(k) > 0 {
			inputBit |= bitCode
		}
	}
	g.bus.UpdateControllerInputBits(0, inputBit) // support single controller for now
}
