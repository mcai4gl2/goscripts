package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	sourcePathPtr := flag.String("source", ".", "source dir containing files to be compressed")
	targetPathPtr := flag.String("target", ".", "target dir to put compressed files in")

	flag.Parse()

	if _, err := os.Stat(*sourcePathPtr); os.IsNotExist(err) {
		panic("source folder doesn't exist")
	}

	var targetPath string = *targetPathPtr

	filepath.Walk(*sourcePathPtr, func(fileFullPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("Skipping folder")
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				os.Mkdir(targetPath, info.Mode())
			}
			return nil
		}

		if strings.HasSuffix(fileFullPath, ".gz") {
			fmt.Println("Already a gz file. Ignoring.")
			return nil
		}

		filename := fileFullPath
		gzedFilename := path.Join(targetPath, info.Name()+".gz")

		if _, err := os.Stat(gzedFilename); os.IsNotExist(err) {
			fmt.Println("Let's create the gz file for " + fileFullPath)
			fmt.Println("New file name: " + gzedFilename)

			originalFile, _ := os.Open(filename)
			defer originalFile.Close()

			reader := bufio.NewReader(originalFile)
			content, _ := ioutil.ReadAll(reader)

			gzFile, err := os.OpenFile(gzedFilename, os.O_CREATE, info.Mode())
			if err != nil {
				panic(err)
			}

			defer gzFile.Close()

			writer := gzip.NewWriter(gzFile)
			defer writer.Close()

			writer.Write(content)
		} else {
			fmt.Println(gzedFilename + " already exists")
		}

		return nil
	})
}
