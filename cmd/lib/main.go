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

	master := internal.MakeMaster(path)

	done := make(chan int)
	worker1 := internal.MakeWorker(1, done)
	worker2 := internal.MakeWorker(2, done)

	start := time.Now()

	go worker1.Run()
	go worker2.Run()

	for i := 0; i < 1; i++ {
		<-done // blocking join operation
	}

	fmt.Printf("Total time: %v\n", time.Since(start).Seconds())
	// get ir files(testing)
	for key := range master.IRfiles {
		fmt.Println(key)
	}
}
