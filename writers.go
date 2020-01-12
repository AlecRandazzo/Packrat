// Copyright (c) 2020 Alec Randazzo

package windowscollector

import (
	"archive/zip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
)

type resultWriter interface {
	ResultWriter(chan fileReader, *sync.WaitGroup) (err error)
}

// ZipResultWriter contains the handles to the file and zip structure
type ZipResultWriter struct {
	ZipWriter  *zip.Writer
	FileHandle *os.File
}

type fileReader struct {
	fullPath string
	reader   io.Reader
}

// ResultWriter will export found files to a zip file.
func (zipResultWriter *ZipResultWriter) ResultWriter(fileReaders chan fileReader, waitForFileCopying *sync.WaitGroup) (err error) {
	defer waitForFileCopying.Done()

	openChannel := true
	for openChannel == true {
		writtenCounter := 0
		fileReader := fileReader{}
		fileReader, openChannel = <-fileReaders
		if openChannel == false {
			break
		}
		normalizedFilePath := strings.ReplaceAll(fileReader.fullPath, "\\", "_")
		normalizedFilePath = strings.ReplaceAll(normalizedFilePath, ":", "_")
		var writer io.Writer
		writer, err = zipResultWriter.ZipWriter.Create(normalizedFilePath)
		if err != nil {
			err = fmt.Errorf("resultWriter failed to add a file to the output zip: %w", err)
			zipResultWriter.ZipWriter.Close()
			zipResultWriter.FileHandle.Close()
			return
		}
		var readErr error
		for {
			buffer := make([]byte, 1024)
			_, readErr = fileReader.reader.Read(buffer)
			if readErr != nil {
				break
			}
			bytesWritten, writeErr := writer.Write(buffer)
			if writeErr != nil {
				log.Panic(writeErr)
			}
			writtenCounter += bytesWritten
		}
		if readErr == io.EOF {
			log.Debugf("Successfully collected '%s'", fileReader.fullPath)
		} else {
			log.Debugf("Failed to collect '%s' due to %v", fileReader.fullPath, readErr)
		}
	}
	zipResultWriter.ZipWriter.Close()
	zipResultWriter.FileHandle.Close()
	err = nil
	return
}
