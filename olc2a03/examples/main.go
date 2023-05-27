package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	gen "nes_emulator/olc2a03/gen"
	"time"
)

func main() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10)) // sr.N(time.Second/10) = buffer size for duration 1/10 second
	done := make(chan bool)

	// Play setting
	const playTime = 2
	wave := gen.NewPulseWave(sr, 220, 0.1)
	wave2 := gen.NewSineWave(sr, 340, 0.3)
	mixer := beep.Mixer{}
	mixer.Add(wave, wave2)
	speaker.Play(beep.Seq(beep.Take(sr.N(playTime*time.Second), &mixer), beep.Callback(func() {
		done <- true
	})))
	<-done
}
