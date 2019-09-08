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
	"fmt"
	mft "github.com/AlecRandazzo/Gofor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
	"syscall"
)

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func (client *CollectorClient) BuildDirectoryTree() (err error) {
	err = client.VolumeHandle.ParseMFTRecord0()
	if err != nil {
		err = fmt.Errorf("failed to parse MFTRecord0: %w", err)
		return
	}
	var waitGroup sync.WaitGroup
	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")

	directoryListChannel := make(chan map[uint64]mft.Directory, 100)
	dataRunsQueue := make(chan mft.DataRun, len(client.VolumeHandle.MftRecord0.DataAttributes.NonResidentDataAttributes.DataRuns))
	for _, value := range client.VolumeHandle.MftRecord0.DataAttributes.NonResidentDataAttributes.DataRuns {
		dataRunsQueue <- value
	}
	close(dataRunsQueue)

	numberOfWorkers := 4
	for i := 0; i < numberOfWorkers; i++ {
		newFileHandle, _ := getVolumeHandle(volumeLetter)
		waitGroup.Add(1)
		go newFileHandle.CreateDirectoryList(&dataRunsQueue, &directoryListChannel, &waitGroup)
	}
	var waitForDirectoryCombination sync.WaitGroup
	waitForDirectoryCombination.Add(1)
	go client.VolumeHandle.CombineDirectoryInformation(&directoryListChannel, &waitForDirectoryCombination)

	waitGroup.Wait()
	close(directoryListChannel)
	waitForDirectoryCombination.Wait()
	log.Debugf("Built trees for %d directories.", len(client.VolumeHandle.MappedDirectories))
	return
}

// Combines a running list of directories from a channel in order to create the systems Directory trees.
func (volume *VolumeHandle) CombineDirectoryInformation(directoryListChannel *chan map[uint64]mft.Directory, waitForDirectoryCombination *sync.WaitGroup) {
	defer waitForDirectoryCombination.Done()

	volume.MappedDirectories = make(map[uint64]string)

	// Merge lists
	var masterDirectoryList map[uint64]mft.Directory
	masterDirectoryList = make(map[uint64]mft.Directory)
	openChannel := true

	for openChannel == true {
		var directoryList map[uint64]mft.Directory
		directoryList = make(map[uint64]mft.Directory)
		directoryList, openChannel = <-*directoryListChannel
		for key, value := range directoryList {
			masterDirectoryList[key] = value
		}
	}

	for recordNumber, directoryMetadata := range masterDirectoryList {
		mappingDirectory := directoryMetadata.DirectoryName
		parentRecordNumberPointer := directoryMetadata.ParentRecordNumber
		for {
			if _, ok := masterDirectoryList[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = ":\\"
					volume.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = ":\\" + mappingDirectory
					volume.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = masterDirectoryList[parentRecordNumberPointer].DirectoryName + "\\" + mappingDirectory
				parentRecordNumberPointer = masterDirectoryList[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			volume.MappedDirectories[recordNumber] = "$ORPHANFILE\\" + mappingDirectory
			break
		}
	}
	return
}

// Creates a list of directories from an MFT read from a volume handle.
func (volume *VolumeHandle) CreateDirectoryList(dataRunQueue *chan mft.DataRun, directoryListChannel *chan map[uint64]mft.Directory, waitGroup *sync.WaitGroup) {
	var directoryList mft.DirectoryList
	var err error
	directoryList = make(map[uint64]mft.Directory)
	openChannel := true

	for openChannel == true {
		dataRun := mft.DataRun{}
		dataRun, openChannel = <-*dataRunQueue
		bytesLeft := dataRun.Length
		_, _ = syscall.Seek(volume.Handle, dataRun.AbsoluteOffset, 0)

		for bytesLeft > 0 {
			var mftRecord mft.MasterFileTableRecord
			buffer := make([]byte, volume.Vbr.MftRecordSize)
			_, _ = syscall.Read(volume.Handle, buffer)
			bytesLeft -= volume.Vbr.MftRecordSize
			mftRecord.MftRecordBytes = buffer

			mftRecord.QuickDirectoryCheck()
			if mftRecord.RecordHeader.FlagDirectory == false {
				continue
			}
			mftRecord.GetRecordHeader()

			err = mftRecord.GetAttributeList()
			if err != nil {
				continue
			}

			err = mftRecord.GetFileNameAttributes()
			if err != nil {
				continue
			}
			for _, attribute := range mftRecord.FileNameAttributes {
				if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") {
					directoryList[uint64(mftRecord.RecordHeader.RecordNumber)] = mft.Directory{DirectoryName: attribute.FileName, ParentRecordNumber: attribute.ParentDirRecordNumber}
					break
				}
			}
		}
	}
	*directoryListChannel <- directoryList
	waitGroup.Done()
	return
}
