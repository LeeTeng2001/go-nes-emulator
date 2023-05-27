package gen

import (
	"github.com/faiface/beep"
	"math"
)

type PulseWave struct {
	baseWave
	dutyCycle float64
}

func (g *PulseWave) SetDutyCycle(perc float64) {
	g.dutyCycle = perc
}

func approxsine(t float64) float64 {
	j := t * 0.15915
	_, j = math.Modf(j)
	return j * 20.785 * (j - 0.5) * (j - 1)
}

// SampleFromSine Fast approx sample from sine
func (g *PulseWave) SampleFromSine(freq, time, duty float64) float64 {
	a, b := 0.0, 0.0
	p := duty * 2 * math.Pi

	for n := 1; n < 20; n++ {
		c := float64(n) * freq * 2 * math.Pi * time
		a += -approxsine(c) / float64(n)
		b += -approxsine(c-p*float64(n)) / float64(n)
	}

	const level = 1
	return (2 * level / math.Pi) * (a - b)
}

func (g *PulseWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if g.time < g.dutyCycle {
			samples[i][0] = g.level * 1
			samples[i][1] = g.level * 1
		} else {
			samples[i][0] = g.level * -1
			samples[i][1] = g.level * -1
		}
		_, g.time = math.Modf(g.time + g.dt)
	}

	return len(samples), true
}

func NewPulseWave(sr beep.SampleRate) *PulseWave {
	s := &PulseWave{}
	s.initBaseWave(sr, 440, 1)
	s.dutyCycle = 0.5 // default
	return s
}
