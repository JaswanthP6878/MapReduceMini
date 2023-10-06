package internal

// global phase
type Phase int

const (
	Map_phase Phase = iota
	Reduce_phase
	Completed_phase
)

// worker phase
type WorkerPhase int

const (
	IDLE WorkerPhase = iota
	Processing
	Done
)
