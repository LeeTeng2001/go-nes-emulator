package main

import (
	"nes_emulator/mlog"
	"nes_emulator/nes"
	"runtime"
)

func main() {
	mlog.L.Infof("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
	g := nes.New()
	g.Run()
}
