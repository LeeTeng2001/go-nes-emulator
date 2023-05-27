package gen

import (
	"github.com/faiface/beep"
	"math"
)

type TriangleWave struct {
	baseWave
}

func (g *TriangleWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = g.level * (math.Abs(math.Mod(2*g.time, 2) - 1))
		samples[i][1] = g.level * (math.Abs(math.Mod(2*g.time, 2) - 1))
		_, g.time = math.Modf(g.time + g.dt)
	}

	return len(samples), true
}

func NewTriangleWave(sr beep.SampleRate, freq float64, level float64) *TriangleWave {
	s := &TriangleWave{}
	s.initBaseWave(sr, freq, level)
	return s
}
