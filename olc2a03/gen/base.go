package gen

import (
	"github.com/faiface/beep"
	"nes_emulator/mlog"
)

type baseWave struct {
	sr          beep.SampleRate
	dt          float64
	time        float64
	level       float64
	streamCount int
}

func (g *baseWave) SetFreq(newFreq float64) {
	dt := newFreq / float64(g.sr)
	if dt >= 0.5 {
		mlog.L.Fatalf("samplerate must be at least 2 times grater then frequency, %v, ratio %v", newFreq, dt)
	}
	g.dt = dt
}

func (g *baseWave) SetLevel(newLevel float64) {
	if newLevel > 1 || newLevel < 0 {
		mlog.L.Fatalf("Level out of range %d", newLevel)
	}
	g.level = newLevel
}

func (*baseWave) Err() error {
	return nil
}

func (g *baseWave) initBaseWave(sr beep.SampleRate, freq float64, level float64) {
	dt := freq / float64(sr)
	if dt >= 1.0/2.0 {
		mlog.L.Fatal("samplerate must be at least 2 times grater then frequency")
	}
	// Protect ear
	if level > 1 {
		mlog.L.Fatalf("Sound level is too high! %v", level)
	}
	g.sr = sr
	g.dt = dt
	g.time = 0
	g.level = level
	g.streamCount = -1 // no limit
}
