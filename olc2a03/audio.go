package olc2a03

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"nes_emulator/bus"
	"nes_emulator/olc2a03/gen"
	"time"
)

// Compile time interface check
var _ bus.ApuDevice = (*Apu)(nil)

type Apu struct {
	// Pulse 1
	pulse1Enable bool
	pulse1Sample float64
	pulse1Seq    sequencer
	pulse1Duty   float64
	// Internal
	frameClockCount int32 // for music timing
	clockCount      int32
	// Audio specification
	playSpeaker *ApuSpeaker
	sampleRate  beep.SampleRate
	// Output buffer for realtime audio synchronisation
	outputBuffer []float64
	bufferClock  int32
	// Generator
	pulseGen *gen.PulseWave
}

// Calculate clock speed of apu because we clock apu together with ppu
// audio is running real time at 44100Hz, 735 samples per frame (60fps)
// Ppu and Apu is clocking at 5369318Hz, 89342 samples per frame
// Approximately apu needs to save output per 121 clocks
const (
	realTimeSampleRate      = 44100
	cpuFreq                 = 1789773
	saveBufferClockInterval = 122
	realTImeSamplePerFrame  = 735
)

func New() *Apu {
	a := &Apu{
		pulse1Enable: false,
		pulse1Sample: 0,
		sampleRate:   beep.SampleRate(realTimeSampleRate),
		outputBuffer: make([]float64, realTImeSamplePerFrame),
		playSpeaker:  NewApuSpeaker(realTImeSamplePerFrame),
		pulseGen:     gen.NewPulseWave(beep.SampleRate(realTimeSampleRate)),
	}
	speaker.Init(a.sampleRate, a.sampleRate.N(time.Second/10))
	speaker.Play(a.playSpeaker)
	return a
}

func (a *Apu) CWrite(addr uint16, data uint8) {
	// https://www.nesdev.org/wiki/APU_basics
	switch addr {
	case 0x4000:
		switch (data & 0xC0) >> 6 { // switch on duty cycle
		case 0x00:
			a.pulse1Seq.sequence = 0b00000001
			a.pulse1Duty = 0.125
		case 0x01:
			a.pulse1Seq.sequence = 0b00000011
			a.pulse1Duty = 0.250
		case 0x02:
			a.pulse1Seq.sequence = 0b00001111
			a.pulse1Duty = 0.500
		case 0x03:
			a.pulse1Seq.sequence = 0b11111100
			a.pulse1Duty = 0.750
		}
	case 0x4001:
	//	pulse1_sweep.enabled = data & 0x80;
	//	pulse1_sweep.period = (data & 0x70) >> 4;
	//	pulse1_sweep.down = data & 0x08;
	//	pulse1_sweep.shift = data & 0x07;
	//	pulse1_sweep.reload = true;
	case 0x4002: // low bit of reload
		a.pulse1Seq.reload = (a.pulse1Seq.reload & 0xFF00) | uint16(data)
	case 0x4003: // high bit of reload and set
		a.pulse1Seq.reload = (a.pulse1Seq.reload & 0x00FF) | (uint16(data&0x7) << 8)
		a.pulse1Seq.timer = a.pulse1Seq.reload
	//	pulse1_seq.sequence = pulse1_seq.new_sequence;
	//	pulse1_lc.counter = length_table[(data & 0xF8) >> 3];
	//	pulse1_env.start = true;
	//	break;
	//
	//case 0x4004:
	//	switch ((data & 0xC0) >> 6)
	//	{
	//	case 0x00: pulse2_seq.new_sequence = 0b01000000; pulse2_osc.dutycycle = 0.125; break;
	//	case 0x01: pulse2_seq.new_sequence = 0b01100000; pulse2_osc.dutycycle = 0.250; break;
	//	case 0x02: pulse2_seq.new_sequence = 0b01111000; pulse2_osc.dutycycle = 0.500; break;
	//	case 0x03: pulse2_seq.new_sequence = 0b10011111; pulse2_osc.dutycycle = 0.750; break;
	//	}
	//	pulse2_seq.sequence = pulse2_seq.new_sequence;
	//	pulse2_halt = (data & 0x20);
	//	pulse2_env.volume = (data & 0x0F);
	//	pulse2_env.disable = (data & 0x10);
	//	break;
	//
	//case 0x4005:
	//	pulse2_sweep.enabled = data & 0x80;
	//	pulse2_sweep.period = (data & 0x70) >> 4;
	//	pulse2_sweep.down = data & 0x08;
	//	pulse2_sweep.shift = data & 0x07;
	//	pulse2_sweep.reload = true;
	//	break;
	//
	//case 0x4006:
	//	pulse2_seq.reload = (pulse2_seq.reload & 0xFF00) | data;
	//	break;
	//
	//case 0x4007:
	//	pulse2_seq.reload = (uint16_t)((data & 0x07)) << 8 | (pulse2_seq.reload & 0x00FF);
	//	pulse2_seq.timer = pulse2_seq.reload;
	//	pulse2_seq.sequence = pulse2_seq.new_sequence;
	//	pulse2_lc.counter = length_table[(data & 0xF8) >> 3];
	//	pulse2_env.start = true;
	//
	//	break;
	//
	//case 0x4008:
	//	break;
	//
	//case 0x400C:
	//	noise_env.volume = (data & 0x0F);
	//	noise_env.disable = (data & 0x10);
	//	noise_halt = (data & 0x20);
	//	break;
	//
	//case 0x400E:
	//	switch (data & 0x0F)
	//	{
	//	case 0x00: noise_seq.reload = 0; break;
	//	case 0x01: noise_seq.reload = 4; break;
	//	case 0x02: noise_seq.reload = 8; break;
	//	case 0x03: noise_seq.reload = 16; break;
	//	case 0x04: noise_seq.reload = 32; break;
	//	case 0x05: noise_seq.reload = 64; break;
	//	case 0x06: noise_seq.reload = 96; break;
	//	case 0x07: noise_seq.reload = 128; break;
	//	case 0x08: noise_seq.reload = 160; break;
	//	case 0x09: noise_seq.reload = 202; break;
	//	case 0x0A: noise_seq.reload = 254; break;
	//	case 0x0B: noise_seq.reload = 380; break;
	//	case 0x0C: noise_seq.reload = 508; break;
	//	case 0x0D: noise_seq.reload = 1016; break;
	//	case 0x0E: noise_seq.reload = 2034; break;
	//	case 0x0F: noise_seq.reload = 4068; break;
	//	}
	//	break;
	//
	case 0x4015: // APU STATUS
		if data&0x01 != 0 {
			a.pulse1Enable = true
		} else {
			a.pulse1Enable = false
		}
		//pulse2_enable = data & 0x02
		//noise_enable = data & 0x04
		//case 0x400F:
		//	pulse1_env.start = true;
		//	pulse2_env.start = true;
		//	noise_env.start = true;
		//	noise_lc.counter = length_table[(data & 0xF8) >> 3];
		//	break;
	}
}

func (a *Apu) CRead(addr uint16) (data uint8) {
	//TODO implement me
	panic("implement me")
}

func (a *Apu) Reset() {
	//TODO implement me
	panic("implement me")
}

func (a *Apu) Clock() {
	quarterFrameClock := false
	halfFrameClock := false

	if a.clockCount%6 == 0 { // Sync to cpu clock
		a.frameClockCount++
		// Magic number! 4 step sequence mode
		if a.frameClockCount == 3729 {
			quarterFrameClock = true
		}
		if a.frameClockCount == 7457 {
			quarterFrameClock = true
			halfFrameClock = true
		}
		if a.frameClockCount == 11186 {
			quarterFrameClock = true
		}
		if a.frameClockCount == 14915 {
			quarterFrameClock = true
			halfFrameClock = true
			a.frameClockCount = 0
		}

		if quarterFrameClock {
			// TODO: Adjust volume envelope
		}
		if halfFrameClock {
			// TODO: Adjust note length and frequency sweeper
		}

		// Update pulse1 channel, rotate will generate us waveform
		a.pulse1Seq.clock(a.pulse1Enable, seqFuncRotate)
		fPulse := cpuFreq / (16 * (float64(a.pulse1Seq.reload) + 1))
		//a.pulse1Sample = a.pulseGen.Sample(fPulse, float64(a.bufferClock)/cpuFreq, a.pulse1Duty)
		a.pulse1Sample = a.pulseGen.SampleFromSine(fPulse, float64(a.bufferClock)/(3*cpuFreq), a.pulse1Duty)
	}
	a.clockCount++

	// Synchronise with real time frame audio buffer
	if a.bufferClock%saveBufferClockInterval == 0 {
		a.outputBuffer[a.bufferClock/saveBufferClockInterval] = a.GetOutputSample()
	}
	a.bufferClock++
}

func (a *Apu) MoveFrameAudioBuffer() {
	a.bufferClock = 0
	a.playSpeaker.UpdateNewBuffer(a.outputBuffer)
	speaker.Lock()
	a.playSpeaker.SwitchBuffer()
	speaker.Unlock()
}

func (a *Apu) GetOutputSample() float64 {
	return a.pulse1Sample
}
