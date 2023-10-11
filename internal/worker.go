package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"sort"
	"time"
)

type Worker struct {
	id           int
	done         chan int // sends done signal to the main function not master
	worker_count int
}

// intermediate files location
const iR_dir string = "/Users/jaswanthpinnepu/Desktop/irfs"
const out_dir string = "/Users/jaswanthpinnepu/Desktop/outfs"

// process input file and returns the IR file location
// sends the IR data files
func (w *Worker) Mapwork(files []string) ([]string, error) {
	var irFiles = []string{}
	for i := 1; i <= w.worker_count; i++ {
		irFiles = append(irFiles, fmt.Sprintf("%v/mr-%v-%v", iR_dir, w.id, i))
	}

	//if IR file does not exist create it
	for _, ir_file := range irFiles {
		if _, err := os.Stat(ir_file); err != nil {
			_, err := os.Create(ir_file)
			if err != nil {
				fmt.Printf("IR file cannot be created!")
				return []string{}, err
			}
		}
	}
	irData := []KeyValue{}
	for _, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("File could not be opened")
			return []string{}, err
		}
		defer file.Close()
		filedata, _ := io.ReadAll(file)
		// call the Map function...
		kv := Map(fileName, string(filedata))
		irData = append(irData, kv...)
	}
	// sort them so that my life becomes easy on reduce phase or not????
	sort.Sort(ByKey(irData))

	// file descriptors
	encoderArray := []json.Encoder{}
	fdarray := []*os.File{}
	for _, irfile := range irFiles {
		fd, err := os.OpenFile(irfile, os.O_RDWR, os.ModeAppend)
		fdarray = append(fdarray, fd)
		if err != nil {
			fmt.Println("Error while opening file")
			return []string{}, err
		}
		encoder := json.NewEncoder(fd)
		encoderArray = append(encoderArray, *encoder)
	}
	// mapping and creating encoders
	// for each of the files

	// enterting values into the encoders
	i := 0
	for i < len(irData) {
		j := i + 1
		for j < len(irData) && irData[j].Key == irData[i].Key {
			j++
		}
		// values := []string{}
		encoderIndex := ihash(irData[i].Key) % w.worker_count
		for k := i; k < j; k++ {
			encoderArray[encoderIndex].Encode(irData[k])
		}
		i = j
	}
	for _, file := range fdarray {
		file.Close()
	}
	return irFiles, nil
}

// reduce function
// files are of the format mr-k-w.id(so have to read all with w.id and parse)
// and apply reduce
func (w *Worker) reduceWork(files []string) (string, error) {
	outFile := fmt.Sprintf("%v/out-%v", out_dir, w.id)
	ofile, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("cannot create outfile\n")
		return "", nil
	}

	irFiles := []string{}
	for _, file := range files {
		// based on the last value of the irFile name
		if string(file[len(file)-1]) == fmt.Sprintf("%v", w.id) {
			irFiles = append(irFiles, file)
		}
	}
	irData := []KeyValue{}
	for _, file := range irFiles {
		fs, err := os.Open(file)
		if err != nil {
			fmt.Printf("reduce task cannot open file")
			return "", err
		}
		dec := json.NewDecoder(fs)
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
				break
			}
			irData = append(irData, kv)
		}
	}

	//sorting
	sort.Sort(ByKey(irData))

	// applying reduce
	i := 0
	for i < len(irData) {
		j := i + 1
		for j < len(irData) && irData[j].Key == irData[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, irData[k].Value)
		}
		output := Reduce(irData[i].Key, values)
		// this is the correct format for each line of Reduce output.
		fmt.Fprintf(ofile, "%v %v\n", irData[i].Key, output)
		i = j
	}
	ofile.Close()
	return outFile, nil
}

// loop of worker lifespan
func (w *Worker) Run() {
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
			fmt.Println("worker_id", w.id, "map phase..., reading input split", reply.TaskType)
			irFiles, err := w.Mapwork(reply.FileName)
			if err != nil {
				fmt.Printf("Error occured in worker, %s", err)
				break
			}
			// send a signal that map work is completed
			args := SetIRfileArgs{IRfiles: irFiles, WorkerId: w.id}
			reply := SetIRFileReply{}
			call("Master.SetIRFile", args, &reply)

		} else if reply.TaskType == Wait {
			time.Sleep(1 * time.Second)

		} else if reply.TaskType == End_phase {
			outFile, err := w.reduceWork(reply.FileName)
			if err != nil {
				fmt.Printf("Reduce phase failed for worker %v", w.id)
				break
			}
			args := EndMRArgs{OutFile: outFile}
			reply := EndMRReply{}
			call("Master.EndMR", args, &reply)
			break
		}
	}
	fmt.Printf("worker %v  has ended!!\n", w.id)
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

func MakeWorker(id int, worker_count int, done chan int) *Worker {
	return &Worker{id: id, done: done, worker_count: worker_count}
}
