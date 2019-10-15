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
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type CollectorClient struct {
	waitGroup               sync.WaitGroup
	VolumeHandler           VolumeHandler
	FileEqualListForFinding fileEqualSearchList
	FileRegexListForFinding fileRegexSearchList
}

type fileReader struct {
	fullPath string
	reader   io.Reader
}

type ResultWriter interface {
	ResultWriter(fileReaders *chan fileReader, waitGroup *sync.WaitGroup, waitForFileCopying *sync.WaitGroup) (err error)
}

type ZipResultWriter struct {
	ZipFileName string
}

// ResultWriter collects target files and writes them to a zip file.
func (zipResultWriter ZipResultWriter) ResultWriter(fileReaders *chan fileReader, waitForInitialization *sync.WaitGroup, waitForFileCopying *sync.WaitGroup) (err error) {
	// Sanity checks
	if zipResultWriter.ZipFileName == "" {
		err = errors.New("ZipResultWriter did not have a ZipFileName set")
		return
	}

	zipFileHandle, err := os.Create(zipResultWriter.ZipFileName)
	if err != nil {
		err = fmt.Errorf("failed to create the zip file %v", zipResultWriter.ZipFileName)
		return
	}
	defer zipFileHandle.Close()
	waitForInitialization.Done()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	var openChannel bool
	writtenCounter := 0
	for {
		reader := fileReader{}
		reader, openChannel = <-*fileReaders
		if openChannel == false {
			break
		}
		log.Debugf("Writing bytes via raw file handle for %s", reader.fullPath)
		writer, _ := zipWriter.Create(reader.fullPath)
		for err != io.EOF {
			buffer := make([]byte, 4096)
			_, err = reader.reader.Read(buffer)
			bytesWritten, _ := writer.Write(buffer)
			writtenCounter += bytesWritten
		}
		log.Debugf("Written a total of %v bytes for the file %v", writtenCounter, reader.fullPath)
	}
	waitForFileCopying.Done()
	err = nil
	return
}
