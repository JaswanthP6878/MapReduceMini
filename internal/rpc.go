package internal

import (
	"os"
	"strconv"
)

// cook-up a unique socket for the system for rpc calls
// to communicate.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
