package gen

import (
	"github.com/faiface/beep"
	"math/rand"
)

type NoiseWave struct {
	baseWave
}

func (g *NoiseWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = g.level * (rand.Float64()*2 - 1)
		samples[i][1] = g.level * (rand.Float64()*2 - 1)
	}

	return len(samples), true
}

func NewNoiseWave(sr beep.SampleRate, level float64) *NoiseWave {
	s := &NoiseWave{}
	s.initBaseWave(sr, 0, level)
	return s
}
