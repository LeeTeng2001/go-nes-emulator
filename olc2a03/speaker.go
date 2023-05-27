package olc2a03

// ApuSpeaker is an infinite speaker listen to pending stream
type ApuSpeaker struct {
	inUseBufferIdx int
	buffers        [2][]float64
	nextSampleIdx  int
	bufferLen      int
}

func NewApuSpeaker(bufferLen int) *ApuSpeaker {
	a := &ApuSpeaker{
		bufferLen: bufferLen,
	}
	for i := 0; i < 2; i++ {
		a.buffers[i] = make([]float64, bufferLen)
	}
	return a
}

// SwitchBuffer should be called with lock, switch cached buffer and reset playhead idx
func (q *ApuSpeaker) SwitchBuffer() {
	q.nextSampleIdx = 0
	if q.inUseBufferIdx == 1 {
		q.inUseBufferIdx = 0
	} else {
		q.inUseBufferIdx = 1
	}
}

// UpdateNewBuffer for single frame
func (q *ApuSpeaker) UpdateNewBuffer(newBuffer []float64) {
	updateIdx := 0
	if q.inUseBufferIdx == 0 {
		updateIdx = 1
	}
	for i := range newBuffer {
		q.buffers[updateIdx][i] = newBuffer[i]
	}
}

func (q *ApuSpeaker) Stream(samples [][2]float64) (int, bool) {
	for i := range samples {
		t := q.buffers[q.inUseBufferIdx][q.nextSampleIdx]
		q.nextSampleIdx = (q.nextSampleIdx + 1) % q.bufferLen
		samples[i][0] = t
		samples[i][1] = t
	}
	return len(samples), true
}

func (q *ApuSpeaker) Err() error {
	return nil
}
