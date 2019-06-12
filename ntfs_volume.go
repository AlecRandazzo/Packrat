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
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"sync"
	"syscall"
)

type VolumeHandle struct {
	Handle            syscall.Handle
	VolumeLetter      string
	Vbr               VolumeBootRecord
	MappedDirectories map[uint64]string
	MftRecord0        mft.MasterFileTableRecord
}

// File that you want to export.
type FileToExport struct {
	FullPath string
	Type     string
}

// Slice of files that you want to export.
type ExportList []FileToExport

type fileEqualListForFinding []string
type fileRegexListForFinding []*regexp.Regexp

type fileExportNameAndBytes struct {
	outputFileName     string
	outputBytesChannel chan []byte
}

// Gets a file handle to the specified volume. This handle is used to read the MFT directly and enables the copying of the MFT despite it being a locked file.
func getVolumeHandle(volumeLetter string) (volume VolumeHandle, err error) {
	const volumeBootRecordSize = 512

	// Get our volume handle
	volume.VolumeLetter = volumeLetter
	volumePath, err := syscall.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s:", volume.VolumeLetter))
	if err != nil {
		err = errors.Wrap(err, "failed to format volume path for syscall")
		return
	}

	dwDesiredAccess := uint32(0x80000000) //0x80 FILE_READ_ATTRIBUTES
	dwShareMode := uint32(0x02 | 0x01)
	dwCreationDisposition := uint32(0x03)
	dwFlagsAndAttributes := uint32(0x00)
	volume.Handle, err = syscall.CreateFile(volumePath, dwDesiredAccess, dwShareMode, nil, dwCreationDisposition, dwFlagsAndAttributes, int32(0))
	if err != nil {
		err = errors.Wrapf(err, "failed to get handle to volume %s", volume.VolumeLetter)
		return
	}

	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, volumeBootRecordSize)
	_, err = syscall.Read(volume.Handle, volumeBootRecord)
	if err != nil {
		err = errors.Wrapf(err, "failed to read %s", volume.VolumeLetter)
		return
	}

	volume.Vbr, err = ParseVolumeBootRecord(volumeBootRecord)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse vbr from volume letter %s", volume.VolumeLetter)
		return
	}
	return
}

func (volume *VolumeHandle) ParseMFTRecord0() (err error) {
	// Move handle pointer back to beginning of volume
	_, err = syscall.Seek(volume.Handle, 0x00, 0)
	if err != nil {
		err = errors.Wrap(err, "failed to see back to volume offset 0x00")
		return
	}

	// Seek to the offset where the MFT starts. If it errors, bomb.
	_, err = syscall.Seek(volume.Handle, volume.Vbr.MftByteOffset, 0)
	if err != nil {
		err = errors.Wrap(err, "failed to seek to mft")
		return
	}

	// Read the first entry in the MFT. The first record in the MFT always is for the MFT itself. If it errors, bomb.
	buffer := make([]byte, volume.Vbr.MftRecordSize)
	_, err = syscall.Read(volume.Handle, buffer)
	if err != nil {
		err = errors.Wrap(err, "failed to read the mft")
		return
	}
	volume.MftRecord0.MftRecordBytes = buffer
	recordHeaderPresent := volume.MftRecord0.CheckForRecordHeader()
	if recordHeaderPresent == false {
		return
	}

	// Everything checks out, let's copy contents of the buffer to the MasterFileRecord struct and then parse the MFT record
	volume.MftRecord0.BytesPerCluster = volume.Vbr.BytesPerCluster
	err = volume.MftRecord0.ParseMFTRecord()
	if err != nil {
		err = errors.Wrap(err, "failed to parse the mft's mft record")
		return
	}
	return
}
func (client *CollectorClient) startCollecting(exportList ExportList) (err error) {
	client.FileEqualListForFinding, client.FileRegexListForFinding, err = buildFileExportLists(exportList)

	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")
	client.VolumeHandle, err = getVolumeHandle(volumeLetter)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Building directory tree.")
	err = client.BuildDirectoryTree()
	if err != nil {
		err = errors.Wrap(err, "Failed to read MFT")
		return
	}
	log.Debugf("Searching the MFT for the following files: %+v", exportList)
	err = client.findFiles()
	if err != nil {
		err = errors.Wrap(err, "failed to findfiles")
		return
	}
	return
}

func (client CollectorClient) mftRecordToBytes(filesToCopyQueue *chan mft.MasterFileTableRecord, fileCopyWaitGroup *sync.WaitGroup) (err error) {
	defer fileCopyWaitGroup.Done()
	openChannel := true

	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")
	volume := VolumeHandle{}
	volume, err = getVolumeHandle(volumeLetter)
	if err != nil {
		log.Fatal(err)
	}

	for openChannel == true {
		var mftRecord mft.MasterFileTableRecord
		mftRecord, openChannel = <-*filesToCopyQueue
		var bytesLeft int64
		var fileName string
		for _, attribute := range mftRecord.FileNameAttributes {
			if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") == true {
				fileName = strings.Replace(client.VolumeHandle.MappedDirectories[attribute.ParentDirRecordNumber]+"_"+attribute.FileName, "\\", "_", -1)
				fileName = strings.Replace(fileName, ":", "_", -1)
				if attribute.LogicalFileSize == 0 || attribute.FileName == "$MFT" {
					for _, value := range mftRecord.DataAttributes.NonResidentDataAttributes.DataRuns {
						bytesLeft = bytesLeft + value.Length
					}
				} else {
					bytesLeft = int64(attribute.LogicalFileSize)
				}
				break
			}
		}

		outputBytesChannel := make(chan []byte, 1000)

		fileExportNameAndBytes := fileExportNameAndBytes{
			outputFileName:     fileName,
			outputBytesChannel: outputBytesChannel,
		}

		client.FileWriteQueue <- fileExportNameAndBytes

		for i := 0; i < len(mftRecord.DataAttributes.NonResidentDataAttributes.DataRuns); i++ {
			offset := mftRecord.DataAttributes.NonResidentDataAttributes.DataRuns[i].AbsoluteOffset
			bytesLeftInRun := mftRecord.DataAttributes.NonResidentDataAttributes.DataRuns[i].Length
			_, _ = syscall.Seek(volume.Handle, offset, 0)

			for bytesLeftInRun > 0 && bytesLeft > 0 {
				if bytesLeft < mftRecord.BytesPerCluster {
					buffer := make([]byte, bytesLeft)
					_, _ = syscall.Read(volume.Handle, buffer)
					outputBytesChannel <- buffer
					break
				} else {
					buffer := make([]byte, mftRecord.BytesPerCluster)
					_, _ = syscall.Read(volume.Handle, buffer)
					outputBytesChannel <- buffer
					bytesLeftInRun -= mftRecord.BytesPerCluster
					bytesLeft -= mftRecord.BytesPerCluster
				}
			}
		}
		close(fileExportNameAndBytes.outputBytesChannel)
	}
	return
}

func buildFileExportLists(exportList ExportList) (fileEqualList fileEqualListForFinding, fileRegexList fileRegexListForFinding, err error) {
	// Normalize everything
	re := regexp.MustCompile("^.:")
	for key, value := range exportList {
		exportList[key].FullPath = strings.ToLower(re.ReplaceAllString(value.FullPath, ":"))
	}

	for _, value := range exportList {
		switch value.Type {
		case "equal":
			fileEqualList = append(fileEqualList, value.FullPath)
		case "regex":
			re := regexp.MustCompile(value.FullPath)
			fileRegexList = append(fileRegexList, re)
		}
	}

	return
}
