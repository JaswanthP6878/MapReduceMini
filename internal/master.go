package internal

import (
	"fmt"
	"os"
)

type Master struct {
	InputFiles []string
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

	return &Master{InputFiles: readfiles}

}
