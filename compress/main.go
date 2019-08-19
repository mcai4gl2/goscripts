package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	sourcePathPtr := flag.String("source", "", "source dir containing files to be compressed")
	targetPathPtr := flag.String("target", ".", "target dir to put compressed files in")
	parallelPtr := flag.Bool("parallel", false, "run compress in parallel")

	flag.Parse()

	if *sourcePathPtr == "" {
		panic("Source dir cannot be empty")
	}

	if _, err := os.Stat(*sourcePathPtr); os.IsNotExist(err) {
		panic("source folder doesn't exist")
	}

	if *parallelPtr {
		log.Println("Running in parallel mode")
	}

	compress(*sourcePathPtr, *targetPathPtr, *parallelPtr)
}
