/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package GoFor_Collector

import (
	"archive/zip"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type CollectorClient struct {
	FileWriteQueue          chan fileExportNameAndBytes
	waitGroup               sync.WaitGroup
	VolumeHandle            VolumeHandle
	FileEqualListForFinding fileEqualListForFinding
	FileRegexListForFinding fileRegexListForFinding
}

// Collects target files and writes them to a zip file.
func (client *CollectorClient) ExportToZip(exportList ExportList, outFileName string) {
	client.FileWriteQueue = make(chan fileExportNameAndBytes)
	client.waitGroup.Add(1)
	go client.startCollecting(exportList)
	client.waitGroup.Wait()
	openChannel := true
	zipFileHandle, _ := os.Create(outFileName)
	defer zipFileHandle.Close()
	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	for openChannel == true {
		fileExportNameAndBytes := fileExportNameAndBytes{}
		fileExportNameAndBytes, openChannel = <-client.FileWriteQueue
		if openChannel == false {
			break
		}
		log.Debugf("Writing bytes for %s", fileExportNameAndBytes.outputFileName)
		fileHandle, _ := zipWriter.Create(fileExportNameAndBytes.outputFileName)
		bytesLeft := true
		for bytesLeft == true {
			var buffer []byte
			buffer, bytesLeft = <-fileExportNameAndBytes.outputBytesChannel
			_, _ = fileHandle.Write(buffer)
		}
		log.Debugf("Finished writing bytes for %s", fileExportNameAndBytes.outputFileName)
	}

	return
}
