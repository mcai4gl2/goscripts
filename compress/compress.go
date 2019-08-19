package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func compressFile(sourceFileFullPath string, targetFileFullPath string, mode os.FileMode) error {
	originalFile, _ := os.Open(sourceFileFullPath)
	defer originalFile.Close()

	reader := bufio.NewReader(originalFile)

	buffer := make([]byte, 1024)

	gzFile, err := os.OpenFile(targetFileFullPath, os.O_CREATE|os.O_WRONLY, mode)

	if err != nil {
		return err
	}

	defer gzFile.Close()

	writer := gzip.NewWriter(gzFile)
	defer writer.Close()

	for {
		num, err := reader.Read(buffer)
		if num != 0 {
			writer.Write(buffer[:num])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			return err
		}
	}

	return nil
}

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
				err = compressFile(fileFullPath, gzedFilename, info.Mode())

				if err != nil {
					panic(err)
				}
			} else {
				waitGroup.Add(1)
				go func(waiter *sync.WaitGroup) {
					err = compressFile(fileFullPath, gzedFilename, info.Mode())

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
