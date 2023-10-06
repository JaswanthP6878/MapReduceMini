package internal

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Master struct {
	InputFiles map[string]bool
	phase      Phase
	sync.Mutex
}

// return filename and phase for the worker
// func (m *Master) AllocateTask() (string, Phase) {
// 	if m.phase == Map_phase {
// 		m.Lock()
// 		defer m.Unlock()
// 		for key, value := range m.InputFiles {
// 			if !value {
// 				m.InputFiles[key] = true // assuming that tasks dont fail at all
// 				return key, m.phase
// 			}
// 		}
// 		fmt.Println("Map phase completed")
// 		return "", 1
// 	}

// }

// testing not actual code.
func (m *Master) AllocateMapTask() (string, Phase) {
	m.Lock()
	defer m.Unlock()
	for key, value := range m.InputFiles {
		if !value {
			m.InputFiles[key] = true // assuming that tasks dont fail at all
			return key, m.phase
		}
	}
	fmt.Println("Map phase completed")
	return "", End_phase

}

// task request by worker
func (m *Master) GetTask(args GetTaskArgs, reply *GetTaskReply) error {
	_ = args.X

	fileName, phase := m.AllocateMapTask()
	reply.FileName = fileName
	reply.TaskType = phase

	return nil

}

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
	mappedFiles := make(map[string]bool)
	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", inputFilesPath, file.Name())
		mappedFiles[filename] = false

	}
	//create the Master with initially in map phase.
	m := Master{InputFiles: mappedFiles, phase: Map_phase}
	m.server()
	return &m

}
