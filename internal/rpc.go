package internal

import (
	"os"
	"strconv"
)

type GetTaskArgs struct {
	X int // see if we need worker id
}

type GetTaskReply struct {
	FileName string
	TaskType Phase
}

// cook-up a unique socket for the system for rpc calls
// to communicate.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
