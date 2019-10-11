/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

// TODO Handle different volumes elegantly

package GoFor_Collector

import (
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
	"os"
	"strings"
	"sync"
)

func (client *CollectorClient) findFilesInDataRun(newVolumeHandle *VolumeHandler, dataRunsQueue *chan mft.DataRun, waitGroup *sync.WaitGroup) {
	client.waitGroup.Done()
	var fileCopyWaitGroup sync.WaitGroup
	openChannel := true

	// Spin up workers that will do the file copying
	numberOfWorkers := 1
	filesToCopyQueue := make(chan mft.MasterFileTableRecord, 100)
	for i := 0; i < numberOfWorkers; i++ {
		fileCopyWaitGroup.Add(1)
		go client.mftRecordToBytes(&filesToCopyQueue, &fileCopyWaitGroup)
	}

	for openChannel == true {
		dataRun := mft.DataRun{}
		dataRun, openChannel = <-*dataRunsQueue
		log.Debugf("Searching for files in the following datarun: %+v", dataRun)
		bytesLeft := dataRun.Length
		_, _ = syscall.Seek(newVolumeHandle.Handle, dataRun.AbsoluteOffset, 0)
		for bytesLeft > 0 {
			buffer := make([]byte, newVolumeHandle.Vbr.MftRecordSize)
			_, _ = syscall.Read(newVolumeHandle.Handle, buffer)

			mftRecord, _ := mft.RawMasterFileTableRecord(buffer).Parse()
			if len(mftRecord.FileNameAttributes) == 0 {
				bytesLeft -= newVolumeHandle.Vbr.MftRecordSize
				continue
			}

			for _, attribute := range mftRecord.FileNameAttributes {
				if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") == true {
					recordFullPath := strings.ToLower(newVolumeHandle.MappedDirectories[attribute.ParentDirRecordNumber] + "\\" + attribute.FileName)
					for _, file := range client.FileEqualListForFinding {
						if file == recordFullPath {
							log.Debugf("Found the MFT record for the file '%s'. Attempting to make a copy of it.", recordFullPath)
							filesToCopyQueue <- mftRecord
							break
						}
					}
					for _, file := range client.FileRegexListForFinding {
						if file.Match([]byte(recordFullPath)) == true {
							log.Debugf("Found the MFT record for the file '%s'. Attempting to make a copy of it.", recordFullPath)
							filesToCopyQueue <- mftRecord
							break
						}
					}
					break
				}
			}
			bytesLeft -= newVolumeHandle.Vbr.MftRecordSize
		}
	}
	close(filesToCopyQueue)
	fileCopyWaitGroup.Wait()
	close(client.FileWriteQueue)
	waitGroup.Done()
	return
}

func (client *CollectorClient) findFiles() (err error) {
	var fileCopyWaitGroup sync.WaitGroup

	dataRunsQueue := make(chan mft.DataRun, len(client.VolumeHandler.MftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns))
	for _, value := range client.VolumeHandler.MftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns {
		dataRunsQueue <- value
	}
	close(dataRunsQueue)

	workerCount := 1
	for i := 0; i < workerCount; i++ {
		volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")
		newVolumeHandle := VolumeHandler{}
		newVolumeHandle, _ = GetVolumeHandle(volumeLetter)
		newVolumeHandle.MappedDirectories = client.VolumeHandler.MappedDirectories
		fileCopyWaitGroup.Add(1)
		go client.findFilesInDataRun(&newVolumeHandle, &dataRunsQueue, &fileCopyWaitGroup)
	}
	fileCopyWaitGroup.Wait()
	return
}
