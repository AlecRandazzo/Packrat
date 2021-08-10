// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// UnResolvedDirectory type is used for creating a directory tree.
type UnResolvedDirectory struct {
	RecordNumber       uint32
	DirectoryName      string
	ParentRecordNumber uint32
}

// UnresolvedDirectoryTree contains a slice of directories that need to be joined to create a UnResolvedDirectory tree.
type UnresolvedDirectoryTree map[uint32]UnResolvedDirectory

// DirectoryTree contains a directory tree.
type DirectoryTree map[uint32]string

// IsThisADirectory will quickly check the bytes of an MFT record to determine if it is a directory or not.
func (rawMftRecord RawMasterFileTableRecord) IsThisADirectory() (result bool, err error) {
	// Sanity checks that the method received good data
	const offsetRecordFlag = 0x16
	const codeDirectory = 0x03
	sizeOfRawMFTRecord := len(rawMftRecord)
	if sizeOfRawMFTRecord == 0 {
		result = false
		err = errors.New("RawMasterFileTableRecord.IsThisADirectory() received nil bytes ")
		return
	}
	if sizeOfRawMFTRecord <= offsetRecordFlag {
		result = false
		err = errors.New("RawMasterFileTableRecord.IsThisADirectory() received not enough bytes ")
		return
	}

	// Skip straight to the offset where the directory flag resides and check if it has the directory flag or not.
	recordFlag := rawMftRecord[offsetRecordFlag]
	if recordFlag == codeDirectory {
		result = true
	} else {
		result = false
	}
	return
}

// ConvertRawMFTRecordToDirectory will take a raw MFT record that is a directory and return the parsed MFT record for it.
func ConvertRawMFTRecordToDirectory(rawMftRecord RawMasterFileTableRecord) (directory UnResolvedDirectory, err error) {
	// Sanity checks that the raw mft record is a directory or not
	result, err := rawMftRecord.IsThisADirectory()
	if result == false {
		err = errors.New("this is not a directory")
		return
	}

	// Get record header bytes
	rawRecordHeader, err := rawMftRecord.GetRawRecordHeader()
	if err != nil {
		err = fmt.Errorf("failed to parse get record header: %w", err)
		return
	}

	// Parse the raw record header
	recordHeader, _ := rawRecordHeader.Parse()

	// Get the raw mft attributes
	rawAttributes, err := rawMftRecord.GetRawAttributes(recordHeader)
	if err != nil {
		err = fmt.Errorf("failed to get raw attributes: %w", err)
		return
	}

	// Find the filename attribute and parse it for its record number, directory name, and parent record number.
	fileNameAttributes, _, _, _, err := rawAttributes.Parse(int64(4096))
	for _, fileNameAttribute := range fileNameAttributes {
		if strings.Contains(fileNameAttribute.FileNamespace, "WIN32") == true || strings.Contains(fileNameAttribute.FileNamespace, "POSIX") {
			directory.RecordNumber = recordHeader.RecordNumber
			directory.DirectoryName = fileNameAttribute.FileName
			directory.ParentRecordNumber = fileNameAttribute.ParentDirRecordNumber
			break
		}
	}
	return
}

// BuildUnresolvedDirectoryTree takes an MFT and does a first pass to find all the directories listed in it. These will form an unresolved UnResolvedDirectory tree that need to be stitched together.
func BuildUnresolvedDirectoryTree(reader io.Reader) (unresolvedDirectoryTree UnresolvedDirectoryTree, err error) {
	unresolvedDirectoryTree = make(UnresolvedDirectoryTree)
	for {
		buffer := make(RawMasterFileTableRecord, 1024)
		_, err = reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}

		directory, err := ConvertRawMFTRecordToDirectory(buffer)
		if err != nil {
			continue
		}
		unresolvedDirectoryTree[directory.RecordNumber] = directory
	}

	return
}

// Resolve combines a running list of directories from a channel in order to create the systems directory trees.
func (unresolvedDirectoryTree UnresolvedDirectoryTree) Resolve(volumeLetter string) (directoryTree DirectoryTree, err error) {
	err = volumeLetterCheck(volumeLetter)
	if err != nil {
		err = fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
		return
	}
	directoryTree = make(DirectoryTree)
	for recordNumber, directoryMetadata := range unresolvedDirectoryTree {
		// Sanity check
		if directoryMetadata.DirectoryName == "" && directoryMetadata.ParentRecordNumber == 0 && directoryMetadata.RecordNumber == 0 {
			continue
		}

		mappingDirectory := directoryMetadata.DirectoryName
		parentRecordNumberPointer := directoryMetadata.ParentRecordNumber
		for {
			if _, ok := unresolvedDirectoryTree[parentRecordNumberPointer]; ok {
				if recordNumber == 5 {
					mappingDirectory = fmt.Sprintf("%s:\\", volumeLetter)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = fmt.Sprintf("%s:\\%s", volumeLetter, mappingDirectory)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = fmt.Sprintf("%s\\%s", unresolvedDirectoryTree[parentRecordNumberPointer].DirectoryName, mappingDirectory)
				parentRecordNumberPointer = unresolvedDirectoryTree[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			directoryTree[recordNumber] = fmt.Sprintf("%s:\\$ORPHANFILE\\%s", volumeLetter, mappingDirectory)
			break
		}
	}
	return
}

// BuildDirectoryTree takes an MFT and creates a directory tree where the slice keys are the mft record number of the UnResolvedDirectory. This record number is importable because files will reference it as its parent mft record number.
func BuildDirectoryTree(reader io.Reader, volumeLetter string) (directoryTree DirectoryTree, err error) {
	err = volumeLetterCheck(volumeLetter)
	if err != nil {
		err = fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
		return
	}
	directoryTree = make(DirectoryTree)
	unresolvedDirectoryTree, _ := BuildUnresolvedDirectoryTree(reader)
	directoryTree, _ = unresolvedDirectoryTree.Resolve(volumeLetter)
	return
}

func volumeLetterCheck(volumeLetter string) (err error) {
	volumeLetterRune := []rune(volumeLetter)
	if volumeLetter == "" {
		err = errors.New("volume letter was blank")
		return
	} else if len(volumeLetterRune) != 1 {
		err = errors.New("volume letter contained more than one character")
		return
	} else if !unicode.IsLetter(volumeLetterRune[0]) {
		err = errors.New("volume letter was not a letter")
		return
	}
	err = nil
	return
}
