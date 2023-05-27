package gen

import (
	"github.com/faiface/beep"
	"math"
)

type SineWave struct {
	baseWave
}

func (g *SineWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		v := g.level * math.Sin(g.time*2.0*math.Pi)
		samples[i][0] = v
		samples[i][1] = v
		_, g.time = math.Modf(g.time + g.dt)
	}

	return len(samples), true
}

func NewSineWave(sr beep.SampleRate, freq float64, level float64) *SineWave {
	s := &SineWave{}
	s.initBaseWave(sr, freq, level)
	return s
}
