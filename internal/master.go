package internal

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Master struct {
	InputFiles []string
	phase      Phase
}

func (m *Master) AllocateTask() {
	if m.phase == Map_phase {
		// allocate a task to worker for map
		// rpc call to a worker
	} else {
		// allocate a reduce task to worker
	}
}

// func (m *Master) (argType T1, replyType *T2) error {

// }

func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func MakeMaster(inputFilesPath string) *Master {
	files, err := os.ReadDir(inputFilesPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	readfiles := []string{}
	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", inputFilesPath, file.Name())
		readfiles = append(readfiles, filename)
	}
	m := Master{InputFiles: readfiles, phase: Map_phase}
	m.server()
	return &m

}
