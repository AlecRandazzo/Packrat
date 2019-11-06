/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

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

type ResultWriter interface {
	ResultWriter(*chan fileReader, *sync.WaitGroup, *sync.WaitGroup) (err error)
}

type ZipResultWriter struct {
	ZipWriter  *zip.Writer
	fileHandle *os.File
}

type fileReader struct {
	fullPath string
	reader   io.Reader
}

func (zipResultWriter *ZipResultWriter) ResultWriter(fileReaders *chan fileReader, waitForInitialization *sync.WaitGroup, waitForFileCopying *sync.WaitGroup) (err error) {
	defer waitForFileCopying.Done()

	openChannel := true
	for openChannel == true {
		writtenCounter := 0
		fileReader := fileReader{}
		fileReader, openChannel = <-*fileReaders
		if openChannel == false {
			break
		}
		normalizedFilePath := strings.ReplaceAll(fileReader.fullPath, "\\", "_")
		normalizedFilePath = strings.ReplaceAll(normalizedFilePath, ":", "_")
		var writer io.Writer
		writer, err = zipResultWriter.ZipWriter.Create(normalizedFilePath)
		if err != nil {
			fmt.Println(err)
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
	zipResultWriter.fileHandle.Close()
	err = nil
	return
}
