package internal

// global phase
type Phase int

const (
	Map_phase Phase = iota
	Merge_phase
	Reduce_phase
	End_phase
)

// worker phase
type WorkerPhase int

const (
	IDLE WorkerPhase = iota
	Processing
	Done
)
