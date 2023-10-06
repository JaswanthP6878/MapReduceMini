package main

import (
	"fmt"

	"mapreduce.jaswantp.com/internal"
)

func main() {
	// lets start with 1 worker and 1 master.
	fmt.Println("Starting Map reduce tasks")

	// assume files are in ~/desktop/dfs
	path := "/Users/jaswanthpinnepu/Desktop/dfs"

	_ = internal.MakeMaster(path)

	done := make(chan int)
	worker1 := internal.MakeWorker(1, done)
	worker2 := internal.MakeWorker(2, done)
	go worker1.Run()
	go worker2.Run()

	for i := 0; i < 2; i++ {
		<-done // blocking join operation
	}

}
