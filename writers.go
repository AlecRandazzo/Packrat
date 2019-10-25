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
	"strings"
	"sync"
)

type ResultWriter interface {
	ResultWriter(*chan fileReader, *sync.WaitGroup, *sync.WaitGroup) (err error)
}

type ZipResultWriter struct {
	ZipWriter zip.Writer
}

type fileReader struct {
	fullPath string
	reader   io.Reader
}

func (zipResultWriter ZipResultWriter) ResultWriter(fileReaders *chan fileReader, waitForInitialization *sync.WaitGroup, waitForFileCopying *sync.WaitGroup) (err error) {
	defer zipResultWriter.ZipWriter.Close()
	defer waitForFileCopying.Done()

	openChannel := true
	for openChannel == true {
		writtenCounter := 0
		fileReader := fileReader{}
		fileReader, openChannel = <-*fileReaders
		if openChannel == false {
			break
		}
		log.Debugf("Reading %s", fileReader.fullPath)
		normalizedFilePath := strings.ReplaceAll(fileReader.fullPath, "\\", "_")
		normalizedFilePath = strings.ReplaceAll(normalizedFilePath, ":", "_")
		var writer io.Writer
		writer, err = zipResultWriter.ZipWriter.Create(normalizedFilePath)
		if err != nil {
			fmt.Println(err)
		}
		var readErr error
		for readErr == nil {
			buffer := make([]byte, 1024)
			_, readErr = fileReader.reader.Read(buffer)
			if readErr != nil {
				log.Debugf("Stopped reading %s due to %v", fileReader.fullPath, err)
				continue
			}
			bytesWritten, writeErr := writer.Write(buffer)
			if writeErr != nil {
				log.Panic(writeErr)
			}
			writtenCounter += bytesWritten
		}
		if readErr == io.EOF {
			log.Debugf("Written a total of %d bytes for the file %s", writtenCounter, fileReader.fullPath)
		}
	}
	err = nil
	return
}
