package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"time"
)

type Worker struct {
	id   int
	done chan int // sends done signal to the main function not master
}

// intermediate files location
const iR_dir string = "/Users/jaswanthpinnepu/Desktop/irfs"

// process input file and returns the IR file location
func (w *Worker) Mapwork(fileName string) (string, error) {
	ir_file := fmt.Sprintf("%v/mr-%v-out", iR_dir, w.id)

	//if IR file does not exist create it
	if _, err := os.Stat(ir_file); err != nil {
		_, err := os.Create(ir_file)
		if err != nil {
			fmt.Printf("IR file cannot be created!")
			return "", err
		}
	}
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("File could not be opened")
		return "", err
	}
	defer file.Close()
	filedata, _ := io.ReadAll(file)

	// call the Map function...
	kv := Map(fileName, string(filedata))

	Irfile, _ := os.OpenFile(ir_file, os.O_RDWR, os.ModeAppend)
	encoder := json.NewEncoder(Irfile)
	for _, value := range kv {
		err = encoder.Encode(value)
		if err != nil {
			fmt.Println("encoding Error")
			return "", err
		}
	}
	defer Irfile.Close()
	return ir_file, nil

}

// loop of worker lifespan
func (w *Worker) Run() {
	var irFileName string
	for {
		args := GetTaskArgs{WorkerID: w.id}
		reply := GetTaskReply{}
		response := call("Master.GetTask", args, &reply)
		if !response {
			fmt.Println("RPC call failed!!")
			break
		}

		if reply.TaskType == Map_phase { // map phase
			var err error
			fmt.Println("worker_id", w.id, "reading..", reply.FileName, reply.TaskType)
			irFileName, err = w.Mapwork(reply.FileName)
			if err != nil {
				fmt.Printf("Error occured in worker, %s", err)
				break
			}
		} else if reply.TaskType == End_phase { // map phase has completed
			args := SetIRfileArgs{FileName: irFileName}
			reply := SetIRFileReply{}
			call("Master.SetIRFile", args, &reply)
			fmt.Printf("sending IR file name to master from worker %v\n", w.id)
			// break

		} else if reply.TaskType == Reduce_phase {
			fmt.Printf("File name for the task is: %v for worker id, %v\n", reply.FileName, w.id)
			break // stopping the worker
		} else if reply.TaskType == Wait {
			time.Sleep(3 * time.Second)
		}

	}
	fmt.Println("worker has ended!!")
	w.done <- 1
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()
	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

func MakeWorker(id int, done chan int) *Worker {
	return &Worker{id: id, done: done}
}
