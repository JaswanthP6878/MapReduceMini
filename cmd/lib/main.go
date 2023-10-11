package main

import (
	"fmt"
	"time"

	"mapreduce.jaswantp.com/internal"
)

func main() {
	// lets start with 1 worker and 1 master.
	fmt.Println("Starting Map reduce tasks")

	// assume files are in ~/desktop/dfs
	path := "/Users/jaswanthpinnepu/Desktop/dfs"

	// worker_count
	var worker_count int = 2
	master := internal.MakeMaster(path, worker_count)
	done := make(chan int)
	workers := []*internal.Worker{}
	for i := 0; i < worker_count; i++ {
		workers = append(workers, internal.MakeWorker(i+1, worker_count, done))
	}
	start := time.Now()
	// worker run
	for _, worker := range workers {
		go worker.Run()
	}
	//  do blocking for workers
	for i := 0; i < worker_count; i++ {
		<-done // blocking join operation
	}
	fmt.Printf("Total time: %v\n", time.Since(start).Seconds())
	fmt.Println("The Output files location:")
	for _, fileName := range master.OutFiles {
		fmt.Println(fileName)
	}
}
