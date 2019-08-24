package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/mcai4gl2/goscripts/compressutil"
)

func compress(sourcePath string, targetPath string, parallel bool) {

	var waitGroup sync.WaitGroup

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

			if !parallel {
				err = compressutil.CompressFile(fileFullPath, gzedFilename, info.Mode())

				if err != nil {
					panic(err)
				}
			} else {
				waitGroup.Add(1)
				go func(waiter *sync.WaitGroup) {
					err = compressutil.CompressFile(fileFullPath, gzedFilename, info.Mode())

					if err != nil {
						panic(err)
					}

					waiter.Done()
				}(&waitGroup)
			}
		} else {
			log.Println(gzedFilename + " already exists")
		}

		return nil
	})

	waitGroup.Wait()
}

func main() {
	sourcePathPtr := flag.String("source", "", "source dir containing files to be compressed")
	targetPathPtr := flag.String("target", ".", "target dir to put compressed files in")
	parallelPtr := flag.Bool("parallel", false, "run compress in parallel with one file per go routine")

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
