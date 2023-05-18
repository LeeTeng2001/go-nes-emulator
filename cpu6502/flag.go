package cpu6502

const (
	FlagC = 1 << iota // carry
	FlagZ             // zero
	FlagI             // disable interrupts
	FlagD             // unused: decimal mode
	FlagB             // break
	FlagU             // unused
	FlagV             // overflow
	FlagN             // negative
)
