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
	"errors"
	"fmt"
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	vbr "github.com/AlecRandazzo/GoFor-VBR-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

type VolumeHandler struct {
	Handle            syscall.Handle
	VolumeLetter      string
	Vbr               vbr.VolumeBootRecord
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

type DataRunReader struct {
	VolumeHandler VolumeHandler
	DataRun       mft.DataRun
}

func getHandle(volumeLetter string) (handle syscall.Handle, err error) {
	dwDesiredAccess := uint32(0x80000000) //0x80 FILE_READ_ATTRIBUTES
	dwShareMode := uint32(0x02 | 0x01)
	dwCreationDisposition := uint32(0x03)
	dwFlagsAndAttributes := uint32(0x00)

	volumePath, _ := syscall.UTF16PtrFromString(fmt.Sprintf("\\\\.\\%s:", volumeLetter))
	handle, err = syscall.CreateFile(volumePath, dwDesiredAccess, dwShareMode, nil, dwCreationDisposition, dwFlagsAndAttributes, 0)
	if err != nil {
		err = fmt.Errorf("getHandle() failed to get handle to volume %s: %w", volumeLetter, err)
		return
	}
	return
}

// Gets a file handle to the specified volume. This handle is used to read the MFT directly and enables the copying of the MFT despite it being a locked file.
func GetVolumeHandle(volumeLetter string) (volume VolumeHandler, err error) {
	const volumeBootRecordSize = 512

	volume.Handle, err = getHandle(volumeLetter)
	if err != nil {
		err = fmt.Errorf("GetVolumeHandle() failed to get handle to volume %s: %w", volume.VolumeLetter, err)
		return
	}

	// Parse the VBR to get details we need about the volume.
	volumeBootRecord := make([]byte, volumeBootRecordSize)
	_, err = syscall.Read(volume.Handle, volumeBootRecord)
	if err != nil {
		err = fmt.Errorf("failed to read %s: %w", volume.VolumeLetter, err)
		return
	}
	volume.Vbr, err = vbr.RawVolumeBootRecord(volumeBootRecord).Parse()
	if err != nil {
		err = fmt.Errorf("GetVolumeHandle() failed to parse vbr from volume letter %s: %w", volume.VolumeLetter, err)
		return
	}
	return
}

func (volume *VolumeHandler) ParseMFTRecord0() (err error) {
	// Move handle pointer back to beginning of volume
	_, err = syscall.Seek(volume.Handle, 0x00, 0)
	if err != nil {
		err = fmt.Errorf("failed to see back to volume offset 0x00: %w", err)
		return
	}

	// Seek to the offset where the MFT starts. If it errors, bomb.
	_, err = syscall.Seek(volume.Handle, volume.Vbr.MftByteOffset, 0)
	if err != nil {
		err = fmt.Errorf("failed to seek to mft: %w", err)
		return
	}

	// Read the first entry in the MFT. The first record in the MFT always is for the MFT itself. If it errors, bomb.
	buffer := make([]byte, volume.Vbr.MftRecordSize)
	_, err = syscall.Read(volume.Handle, buffer)
	if err != nil {
		err = fmt.Errorf("failed to read the mft: %w", err)
		return
	}

	// Sanity check that this is indeed an mft record
	result, err := mft.RawMasterFileTableRecord(buffer).IsThisAnMftRecord()
	if err != nil {
		err = fmt.Errorf("IsThisAnMftRecord() returned an error: %v", err)
	} else if result == false {
		err = errors.New("VolumeHandler.ParseMFTRecord0() received an invalid mft record")
		return
	}

	// Parse the MFT record

	volume.MftRecord0, err = mft.RawMasterFileTableRecord(buffer).Parse(volume.Vbr.BytesPerCluster)
	if err != nil {
		err = fmt.Errorf("VolumeHandler.ParseMFTRecord0() failed to parse the mft's mft record: %w", err)
		return
	}
	return
}

func (dataRunReader DataRunReader) Read(byteSliceToPopulate []byte) (bytesWritten int, err error) {
	// Get current offset
	currentOffset, err := syscall.Seek(dataRunReader.VolumeHandler.Handle, 0, 1)
	if err != nil {
		log.Println(err)
	}

	// Calculate range current offset needs to be in
	dataRunMin := dataRunReader.DataRun.AbsoluteOffset
	dataRunMax := dataRunReader.DataRun.AbsoluteOffset + (dataRunReader.DataRun.Length * dataRunReader.VolumeHandler.Vbr.BytesPerCluster)

	// Check to see if the offset is outside of the range
	if currentOffset < dataRunMin || currentOffset > dataRunMax {
		// Seek to the beginning of the data run
		_, err = syscall.Seek(dataRunReader.VolumeHandler.Handle, dataRunReader.DataRun.AbsoluteOffset, 0)
		if err != nil {
			log.Println(err)
		}
	}

	// Read the data
	bytesWritten, err = syscall.Read(dataRunReader.VolumeHandler.Handle, byteSliceToPopulate)
	if err != nil {
		log.Println(err)
	}

	// Get the new current offset
	currentOffset, err = syscall.Seek(dataRunReader.VolumeHandler.Handle, 0, 1)
	if err != nil {
		log.Println(err)
	}

	// Check if the new offset is now beyond the scope of the data run. If it is set err to io.EOF
	if currentOffset > dataRunMax {
		err = io.EOF
	}

	return
}

func (client *CollectorClient) startCollecting(exportList ExportList) (err error) {
	client.FileEqualListForFinding, client.FileRegexListForFinding, err = buildFileExportLists(exportList)

	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")
	client.VolumeHandler, err = GetVolumeHandle(volumeLetter)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Building directory tree.")
	unresolvedDirectoryTree := mft.UnresolvedDirectoryTree{}
	for _, dataRun := range client.VolumeHandler.MftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns {
		dataRunReader := DataRunReader{
			VolumeHandler: client.VolumeHandler,
			DataRun:       dataRun,
		}
		tempUnresolvedDirectoryTree, err := mft.BuildUnresolvedDirectoryTree(dataRunReader)
		if err != nil {
			log.Println(err)
			return
		}

		// Merge temporary directory tree with the master tree
		for recordNumber, directory := range tempUnresolvedDirectoryTree {
			unresolvedDirectoryTree[recordNumber] = directory
		}
	}
	client.VolumeHandler.MappedDirectories = unresolvedDirectoryTree.Resolve()

	log.Debugf("Searching the MFT for the following files: %+v", exportList)
	err = client.findFiles()
	if err != nil {
		err = fmt.Errorf("failed to findfiles: %w", err)
		return
	}
	return
}

func (client CollectorClient) mftRecordToBytes(filesToCopyQueue *chan mft.MasterFileTableRecord, fileCopyWaitGroup *sync.WaitGroup) (err error) {
	defer fileCopyWaitGroup.Done()
	openChannel := true

	volumeLetter := strings.TrimRight(os.Getenv("SYSTEMDRIVE"), ":")
	volume := VolumeHandler{}
	volume, err = GetVolumeHandle(volumeLetter)
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
				fileName = strings.Replace(client.VolumeHandler.MappedDirectories[attribute.ParentDirRecordNumber]+"_"+attribute.FileName, "\\", "_", -1)
				fileName = strings.Replace(fileName, ":", "_", -1)
				if attribute.LogicalFileSize == 0 || attribute.FileName == "$MFT" {
					for _, value := range mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns {
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

		for i := 0; i < len(mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns); i++ {
			offset := mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns[i].AbsoluteOffset
			bytesLeftInRun := mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns[i].Length
			_, _ = syscall.Seek(volume.Handle, offset, 0)

			for bytesLeftInRun > 0 && bytesLeft > 0 {
				if bytesLeft < client.VolumeHandler.Vbr.BytesPerCluster {
					buffer := make([]byte, bytesLeft)
					_, _ = syscall.Read(volume.Handle, buffer)
					outputBytesChannel <- buffer
					break
				} else {
					buffer := make([]byte, client.VolumeHandler.Vbr.BytesPerCluster)
					_, _ = syscall.Read(volume.Handle, buffer)
					outputBytesChannel <- buffer
					bytesLeftInRun -= client.VolumeHandler.Vbr.BytesPerCluster
					bytesLeft -= client.VolumeHandler.Vbr.BytesPerCluster
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
