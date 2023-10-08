package internal

import (
	"os"
	"strconv"
)

type GetTaskArgs struct {
	WorkerID int // see if we need worker id
}

type GetTaskReply struct {
	FileName string
	TaskType Phase
}

// setting IR files to master
type SetIRfileArgs struct {
	FileName string
}

type SetIRFileReply struct {
	Ok int
}

// cook-up a unique socket for the system for rpc calls
// to communicate.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
