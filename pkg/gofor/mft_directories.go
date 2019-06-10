/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
)

type directory struct {
	DirectoryName      string
	ParentRecordNumber uint64
}

type directoryList map[uint64]directory

type mappedDirectories map[uint64]string

func (mftRecord *masterFileTableRecord) quickDirectoryCheck() {

	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	if len(mftRecord.MftRecordBytes) <= offsetRecordFlag {
		mftRecord.RecordHeader.FlagDirectory = false
		return
	}
	recordFlag := mftRecord.MftRecordBytes[offsetRecordFlag]
	if recordFlag == codeDirectory {
		mftRecord.RecordHeader.FlagDirectory = true
	} else {
		mftRecord.RecordHeader.FlagDirectory = false
	}
	return
}

func (volume *volumeHandle) createDirectoryList(dataRunQueue *chan dataRun, directoryListChannel *chan map[uint64]directory, waitGroup *sync.WaitGroup) {
	var directoryList directoryList
	var err error
	directoryList = make(map[uint64]directory)
	openChannel := true

	for openChannel == true {
		dataRun := dataRun{}
		dataRun, openChannel = <-*dataRunQueue
		bytesLeft := dataRun.Length
		_, _ = syscall.Seek(volume.Handle, dataRun.AbsoluteOffset, 0)

		for bytesLeft > 0 {
			var mftRecord masterFileTableRecord
			buffer := make([]byte, volume.Vbr.MftRecordSize)
			_, _ = syscall.Read(volume.Handle, buffer)
			bytesLeft -= volume.Vbr.MftRecordSize
			mftRecord.MftRecordBytes = buffer

			mftRecord.quickDirectoryCheck()
			if mftRecord.RecordHeader.FlagDirectory == false {
				continue
			}
			mftRecord.getRecordHeader()

			err = mftRecord.getAttributeList()
			if err != nil {
				continue
			}

			err = mftRecord.getFileNameAttributes()
			if err != nil {
				continue
			}
			for _, attribute := range mftRecord.FileNameAttributes {
				if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") {
					directoryList[uint64(mftRecord.RecordHeader.RecordNumber)] = directory{DirectoryName: attribute.FileName, ParentRecordNumber: attribute.ParentDirRecordNumber}
					break
				}
			}
		}
	}
	*directoryListChannel <- directoryList
	waitGroup.Done()
	return
}

func createDirectoryList(inboundBuffer *chan []byte, directoryListChannel *chan map[uint64]directory, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	var openChannel = true
	var directoryList directoryList
	var err error
	directoryList = make(map[uint64]directory)
	for openChannel == true {
		var mftRecord masterFileTableRecord
		mftRecord.MftRecordBytes, openChannel = <-*inboundBuffer
		mftRecord.quickDirectoryCheck()
		if mftRecord.RecordHeader.FlagDirectory == false {
			continue
		}
		mftRecord.getRecordHeader()

		err = mftRecord.getAttributeList()
		if err != nil {
			continue
		}

		err = mftRecord.getFileNameAttributes()
		if err != nil {
			continue
		}
		for _, attribute := range mftRecord.FileNameAttributes {
			if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") {

				directoryList[uint64(mftRecord.RecordHeader.RecordNumber)] = directory{DirectoryName: attribute.FileName, ParentRecordNumber: attribute.ParentDirRecordNumber}
				break
			}
		}
	}
	*directoryListChannel <- directoryList
	return
}

func (file *mftFile) combineDirectoryInformation(directoryListChannel *chan map[uint64]directory, waitForDirectoryCombination *sync.WaitGroup) {
	defer waitForDirectoryCombination.Done()

	file.MappedDirectories = make(map[uint64]string)

	// Merge lists
	var masterDirectoryList map[uint64]directory
	masterDirectoryList = make(map[uint64]directory)
	openChannel := true

	for openChannel == true {
		var directoryList map[uint64]directory
		directoryList = make(map[uint64]directory)
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
					file.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = ":\\" + mappingDirectory
					file.MappedDirectories[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = masterDirectoryList[parentRecordNumberPointer].DirectoryName + "\\" + mappingDirectory
				parentRecordNumberPointer = masterDirectoryList[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			file.MappedDirectories[recordNumber] = "$ORPHANFILE\\" + mappingDirectory
			break
		}
	}
	return
}

func (volume *volumeHandle) combineDirectoryInformation(directoryListChannel *chan map[uint64]directory, waitForDirectoryCombination *sync.WaitGroup) {
	defer waitForDirectoryCombination.Done()

	volume.MappedDirectories = make(map[uint64]string)

	// Merge lists
	var masterDirectoryList map[uint64]directory
	masterDirectoryList = make(map[uint64]directory)
	openChannel := true

	for openChannel == true {
		var directoryList map[uint64]directory
		directoryList = make(map[uint64]directory)
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

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func (file *mftFile) buildDirectoryTree() (err error) {
	var waitGroup sync.WaitGroup
	bufferChannel := make(chan []byte, 100)
	numberOfWorkers := 4
	directoryListChannel := make(chan map[uint64]directory, numberOfWorkers)
	for i := 0; i <= numberOfWorkers; i++ {
		waitGroup.Add(1)
		go createDirectoryList(&bufferChannel, &directoryListChannel, &waitGroup)
	}

	var waitForDirectoryCombination sync.WaitGroup
	waitForDirectoryCombination.Add(1)
	go file.combineDirectoryInformation(&directoryListChannel, &waitForDirectoryCombination)
	var offset int64 = 0
	for {
		buffer := make([]byte, 1024)
		_, err = file.FileHandle.ReadAt(buffer, offset)
		if err == io.EOF {
			err = nil
			break
		}
		bufferChannel <- buffer
		offset += 1024
	}

	close(bufferChannel)
	waitGroup.Wait()
	close(directoryListChannel)
	waitForDirectoryCombination.Wait()
	return
}

// Builds a list of directories for the purpose of of mapping MFT records to their parent directories.
func (client *CollectorClient) buildDirectoryTree() (err error) {
	err = client.VolumeHandle.parseMFTRecord0()
	if err != nil {
		err = errors.Wrap(err, "Failed to parse MFTRecord0")
		return
	}
	var waitGroup sync.WaitGroup
	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")

	directoryListChannel := make(chan map[uint64]directory, 100)
	dataRunsQueue := make(chan dataRun, len(client.VolumeHandle.MftRecord0.DataAttributes.NonResidentDataAttributes.DataRuns))
	for _, value := range client.VolumeHandle.MftRecord0.DataAttributes.NonResidentDataAttributes.DataRuns {
		dataRunsQueue <- value
	}
	close(dataRunsQueue)

	numberOfWorkers := 4
	for i := 0; i < numberOfWorkers; i++ {
		newFileHandle, _ := getVolumeHandle(volumeLetter)
		waitGroup.Add(1)
		go newFileHandle.createDirectoryList(&dataRunsQueue, &directoryListChannel, &waitGroup)
	}
	var waitForDirectoryCombination sync.WaitGroup
	waitForDirectoryCombination.Add(1)
	go client.VolumeHandle.combineDirectoryInformation(&directoryListChannel, &waitForDirectoryCombination)

	waitGroup.Wait()
	close(directoryListChannel)
	waitForDirectoryCombination.Wait()
	log.Debugf("Built trees for %d directories.", len(client.VolumeHandle.MappedDirectories))
	return
}
