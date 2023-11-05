package internal

// global phase
type Phase int

const (
	Map_phase Phase = iota
	Reduce_phase
	End_phase
	Wait
)

// worker phase
type WorkerPhase int

const (
	Start WorkerPhase = iota
	IDLE
	Processing
	Done
)
