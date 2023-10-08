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
	InputFiles   map[string]bool
	workerFiles  map[int][]string
	phase        Phase
	IRfiles      map[string]bool
	Worker_Count int
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
	for key, value := range m.InputFiles {
		if !value {
			m.InputFiles[key] = true // assuming that tasks dont fail at all
			return key, Map_phase
		}
	}
	// need to write better logic here have to see why its failing
	m.Worker_Count -= 1
	if m.Worker_Count == 0 {
		m.phase = End_phase
		fmt.Println("Reached end for the map phase for all workers")
		return "", End_phase
	} else if m.Worker_Count != 0 {
		return "", Wait
	}
	return "", End_phase

}

func (m *Master) AllocateReduceTask(workerID int) (string, Phase) {
	filePath := fmt.Sprintf("mr-%v-out", workerID)
	return filePath, Reduce_phase
}

// task request by worker
func (m *Master) GetTask(args GetTaskArgs, reply *GetTaskReply) error {
	m.Lock()
	defer m.Unlock()
	workerId := args.WorkerID
	var fileName string
	var phase Phase
	if m.phase == Map_phase {
		fileName, phase = m.AllocateMapTask()
	} else if m.phase == End_phase {
		fileName, phase = m.AllocateReduceTask(workerId)
	} else {
		reply.FileName = fileName
		reply.TaskType = Wait
		return nil
	}
	reply.FileName = fileName
	reply.TaskType = phase
	return nil
}

func (m *Master) SetValue(fileName string) error {
	m.Lock()
	defer m.Unlock()
	m.IRfiles[fileName] = false // adding IR files
	return nil
}

// set The IR file from  the worker:
func (m *Master) SetIRFile(args SetIRfileArgs, reply *SetIRFileReply) error {
	m.SetValue(args.FileName)
	reply.Ok = 1
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

func MakeMaster(inputFilesPath string, workerCount int) *Master {
	files, err := os.ReadDir(inputFilesPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	workerFiles := splitInputDir(inputFilesPath, workerCount)
	for key, files := range workerFiles {
		fmt.Println(key)
		for _, file := range files {
			fmt.Println(file)
		}
	}
	mappedFiles := make(map[string]bool)
	irFiles := make(map[string]bool)
	for _, file := range files {
		filename := fmt.Sprintf("%s/%s", inputFilesPath, file.Name())
		mappedFiles[filename] = false
	}
	//create the Master with initially in map phase.
	m := Master{InputFiles: mappedFiles, phase: Map_phase, IRfiles: irFiles, Worker_Count: workerCount, workerFiles: workerFiles}
	m.server()
	return &m

}
