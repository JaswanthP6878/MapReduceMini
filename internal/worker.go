package internal

import (
	"fmt"
	"log"
	"net/rpc"
)

type Worker struct {
	id   int
	done chan int
}

func (w *Worker) Run() {
	for {
		args := GetTaskArgs{X: 1}
		reply := GetTaskReply{}
		response := call("Master.GetTask", args, &reply)
		if response {
			if reply.TaskType != 1 {

				fmt.Println("worker_id", w.id, "reading..", reply.FileName, reply.TaskType)
			} else {
				fmt.Println("map files have been completed")
				break
			}
		} else {
			fmt.Println("RPC call failed!!")
			break
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
