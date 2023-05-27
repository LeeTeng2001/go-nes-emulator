package olc2a03

type sequencer struct {
	sequence uint32 // some data stored in sequencer
	timer    uint16
	reload   uint16 // what the counter get reset to
	output   uint8
}

// Run some lambda function when the timer reaches 0
func (s *sequencer) clock(enabled bool, funcMap func(seq *uint32)) {
	if enabled {
		s.timer--
		if s.timer == 0xFFFF { // sensitive to -1
			s.timer = s.reload + 1
			funcMap(&s.sequence)
			s.output = uint8(s.sequence & 1)
		}
	}
}

func seqFuncRotate(seq *uint32) {
	// shift right 1 bit
	*seq = ((*seq & 1) << 7) | ((*seq & 0xFE) >> 1)
}
