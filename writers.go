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

type ResultWriter interface {
	ResultWriter(*chan fileReader, *sync.WaitGroup, *sync.WaitGroup) (err error)
}

type ZipResultWriter struct {
	ZipFileName string
}

type fileReader struct {
	fullPath string
	reader   io.Reader
}

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

	for openChannel == true {
		writtenCounter := 0
		fileReader := fileReader{}
		fileReader, openChannel = <-*fileReaders
		if openChannel == false {
			break
		}
		log.Debugf("Reading %s", fileReader.fullPath)
		writer, err := zipWriter.Create(fileReader.fullPath)
		if err != nil {
			fmt.Println(err)
		}
		var readErr error
		for readErr == nil {
			buffer := make([]byte, 4096)
			_, readErr = fileReader.reader.Read(buffer)
			if readErr != nil {
				continue
			}
			bytesWritten, writeErr := writer.Write(buffer)
			if writeErr != nil {
				fmt.Println(writeErr)
			}
			writtenCounter += bytesWritten
		}
		if readErr == io.EOF {
			log.Debugf("Written a total of %d bytes for the file %s", writtenCounter, fileReader.fullPath)
		}
	}
	waitForFileCopying.Done()
	err = nil
	return
}
