package compressutil

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
)

func CompressFile(sourceFileFullPath string, targetFileFullPath string, mode os.FileMode) error {
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
