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

	Master := internal.MakeMaster(path)

	for _, name := range Master.InputFiles {
		fmt.Println(name)
	}

}
