package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mcai4gl2/goscripts/compressutil"
)

type CompressWork struct {
	sourceFileFullPath string
	targetFileFullPath string
	mode               os.FileMode
}

func compressWorker(id int, works <-chan CompressWork) <-chan int {
	done := make(chan int)
	go func() {
		counter := 0
		for work := range works {
			log.Println(fmt.Sprintf("Worker %d: starting to compress %s",
				id, work.sourceFileFullPath))
			compressutil.CompressFile(work.sourceFileFullPath,
				work.targetFileFullPath,
				work.mode)
			log.Println("done")
			counter++
		}
		done <- counter
		close(done)
	}()
	return done
}

func pushWorks(works chan<- CompressWork, sourcePath string, targetPath string) {
	filepath.Walk(sourcePath, func(fileFullPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Println("Skipping folder")
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				os.Mkdir(targetPath, info.Mode())
			}
			return nil
		}

		if strings.HasSuffix(fileFullPath, ".gz") {
			log.Println("Already a gz file. Ignoring.")
			return nil
		}

		gzedFilename := path.Join(targetPath, info.Name()+".gz")

		if _, err := os.Stat(gzedFilename); os.IsNotExist(err) {
			log.Println("Let's create the gz file for " + fileFullPath)
			log.Println("New file name: " + gzedFilename)

			works <- CompressWork{fileFullPath, gzedFilename, info.Mode()}
		} else {
			log.Println(gzedFilename + " already exists")
		}

		return nil
	})

	close(works)
}

func compress(sourcePath string, targetPath string, parallel int) {
	workChannel := make(chan CompressWork)

	workers := make([]<-chan int, parallel)
	for i := 0; i < parallel; i++ {
		workers[i] = compressWorker(i, workChannel)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(workers))

	for index, ch := range workers {
		go func(index int, channel <-chan int) {
			defer waitGroup.Done()
			for i := range channel {
				log.Println(fmt.Sprintf("Worker %d has compressed %d num of files",
					index, i))
			}
		}(index, ch)
	}

	pushWorks(workChannel, sourcePath, targetPath)

	waitGroup.Wait()
}

func main() {
	sourcePathPtr := flag.String("source", "", "source dir containing files to be compressed")
	targetPathPtr := flag.String("target", ".", "target dir to put compressed files in")
	parallelPtr := flag.Int("parallel", 10, "max number of concurrent go routine to compress")

	flag.Parse()

	if *sourcePathPtr == "" {
		panic("Source dir cannot be empty")
	}

	if _, err := os.Stat(*sourcePathPtr); os.IsNotExist(err) {
		panic("source folder doesn't exist")
	}

	log.Println(fmt.Sprintf("Running with %d degree of parallelism",
		*parallelPtr))

	compress(*sourcePathPtr, *targetPathPtr, *parallelPtr)
}
