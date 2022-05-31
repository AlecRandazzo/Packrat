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

const offsetRecordFlag = 0x16

// ValidateDirectory will check the input of an MFT record to determine if it is a directory or not.
func ValidateDirectory(input []byte) error {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return errors.New("received nil input ")
	}

	if size <= offsetRecordFlag {
		return errors.New("received not enough input ")
	}

	// Skip straight to the offset where the directory flag resides and check if it has the directory flag or not.
	recordFlag := input[offsetRecordFlag]
	if recordFlag != codeDirectory {
		return errors.New("mft record is not a directory")
	}

	return nil
}

// ConvertRawMFTRecordToDirectory will take a raw MFT record that is a directory and return the parsed MFT record for it.
func ConvertRawMFTRecordToDirectory(input []byte, bytesPerSector uint) (UnResolvedDirectory, error) {
	// Sanity checks that the raw mft record is a directory or not
	err := ValidateDirectory(input)
	if err != nil {
		return UnResolvedDirectory{}, errors.New("this is not a directory")
	}

	err = fixup(input, bytesPerSector)
	if err != nil {
		return UnResolvedDirectory{}, fmt.Errorf("failed to fixup record: %w", err)
	}

	// Parse the raw record header
	var recordHeader RecordHeader
	recordHeader, err = GetRecordHeaders(input)
	if err != nil {
		return UnResolvedDirectory{}, fmt.Errorf("failed to get record header: %w", err)
	}

	// get raw attributes
	var rawAttributes [][]byte
	rawAttributes, err = GetRawAttributes(input, recordHeader)
	if err != nil {
		return UnResolvedDirectory{}, fmt.Errorf("failed to get raw attributes: %w", err)
	}

	// init return variable
	var directory UnResolvedDirectory

	// Find the filename attribute and parse it for its record number, directory name, and parent record number.
	var fnAttributes FileNameAttributes
	fnAttributes, _, _, _, err = GetAttributes(rawAttributes, 4096)
	for _, fileNameAttribute := range fnAttributes {
		if strings.Contains(string(fileNameAttribute.FileNamespace), "WIN32") ||
			strings.Contains(string(fileNameAttribute.FileNamespace), "WIN32 & DOS") ||
			strings.Contains(string(fileNameAttribute.FileNamespace), "POSIX") {
			directory.RecordNumber = recordHeader.RecordNumber
			directory.DirectoryName = fileNameAttribute.FileName
			directory.ParentRecordNumber = fileNameAttribute.ParentDirRecordNumber
			break
		}
	}
	return directory, nil
}

// buildUnresolvedDirectoryTree takes an MFT and does a first pass to find all the directories listed in it. These will form an unresolved UnResolvedDirectory tree that need to be stitched together.
func buildUnresolvedDirectoryTree(reader io.Reader, bytesPerSector uint) (UnresolvedDirectoryTree, error) {
	tree := make(UnresolvedDirectoryTree)
	for {
		buffer := make([]byte, 1024)
		_, err := reader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}

		var directory UnResolvedDirectory
		directory, err = ConvertRawMFTRecordToDirectory(buffer, bytesPerSector)
		if err != nil {
			continue
		}
		tree[directory.RecordNumber] = directory
	}

	return tree, nil
}

// Resolve combines a running list of directories from a channel in order to create the systems directory trees.
func (unresolvedDirectoryTree UnresolvedDirectoryTree) Resolve(volumeLetter string) (DirectoryTree, error) {
	err := checkVolumeLetter(volumeLetter)
	if err != nil {
		return DirectoryTree{}, fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
	}
	directoryTree := make(DirectoryTree)
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
					mappingDirectory = fmt.Sprintf(`%s:`, volumeLetter)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				if parentRecordNumberPointer == 5 {
					mappingDirectory = fmt.Sprintf(`%s:\%s`, volumeLetter, mappingDirectory)
					directoryTree[recordNumber] = mappingDirectory
					break
				}
				mappingDirectory = fmt.Sprintf(`%s\%s`, unresolvedDirectoryTree[parentRecordNumberPointer].DirectoryName, mappingDirectory)
				parentRecordNumberPointer = unresolvedDirectoryTree[parentRecordNumberPointer].ParentRecordNumber
				continue
			}
			directoryTree[recordNumber] = fmt.Sprintf(`%s:\$ORPHANFILE\%s`, volumeLetter, mappingDirectory)
			break
		}
	}
	return directoryTree, nil
}

// BuildDirectoryTree takes an MFT and creates a directory tree where the slice keys are the mft record number of the UnResolvedDirectory.
func BuildDirectoryTree(reader io.Reader, volumeLetter string, bytesPerSector uint) (DirectoryTree, error) {
	err := checkVolumeLetter(volumeLetter)
	if err != nil {
		return DirectoryTree{}, fmt.Errorf("failed to build directory tree due to invalid volume letter: %w", err)
	}
	directoryTree := make(DirectoryTree)
	unresolvedDirectoryTree, _ := buildUnresolvedDirectoryTree(reader, bytesPerSector)
	directoryTree, _ = unresolvedDirectoryTree.Resolve(volumeLetter)
	return directoryTree, nil
}

func checkVolumeLetter(volumeLetter string) error {
	volumeLetterRune := []rune(volumeLetter)
	if volumeLetter == "" {
		return errors.New("volume letter was blank")
	} else if len(volumeLetterRune) != 1 {
		return errors.New("volume letter contained more than one character")
	} else if !unicode.IsLetter(volumeLetterRune[0]) {
		return errors.New("volume letter was not a letter")
	}
	return nil
}
