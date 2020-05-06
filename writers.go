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

	// We receive io.Readers from the fileReaders channel. These are files that the collector identified as ones to collect.
	openChannel := true
	for openChannel == true {
		fileReader := fileReader{}
		fileReader, openChannel = <-fileReaders
		if openChannel == false {
			break
		}

		// Normalize the file path so we can make the path a valid file name
		normalizedFilePath := strings.ReplaceAll(fileReader.fullPath, "\\", "_")
		normalizedFilePath = strings.ReplaceAll(normalizedFilePath, ":", "_")

		// Create a new file inside the zip file
		var writer io.Writer
		writer, err = zipResultWriter.ZipWriter.Create(normalizedFilePath)
		if err != nil {
			err = fmt.Errorf("resultWriter failed to add a file to the output zip: %w", err)
			return
		}

		// Copy contents from the file we want to collect to the output file in the zip
		_, readErr := io.Copy(writer, fileReader.reader)
		if readErr == io.EOF {
			log.Debugf("Successfully collected '%s'", fileReader.fullPath)
		} else {
			log.Debugf("Failed to collect '%s' due to %v", fileReader.fullPath, readErr)
		}
	}
	err = nil
	return
}
