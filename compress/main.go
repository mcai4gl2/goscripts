package main

import (
	"flag"
	"os"
)

func main() {
	sourcePathPtr := flag.String("source", "", "source dir containing files to be compressed")
	targetPathPtr := flag.String("target", ".", "target dir to put compressed files in")

	flag.Parse()

	if *sourcePathPtr == "" {
		panic("Source dir cannot be empty")
	}

	if _, err := os.Stat(*sourcePathPtr); os.IsNotExist(err) {
		panic("source folder doesn't exist")
	}

	var targetPath string = *targetPathPtr

	compress(*sourcePathPtr, targetPath)
}
