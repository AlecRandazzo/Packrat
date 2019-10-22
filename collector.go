package windowscollector

import (
	"fmt"
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	syscall "golang.org/x/sys/windows"
	"io"
	"strings"
	"sync"
)

func Collect(volumeLetter string, exportList ListOfFilesToExport, resultWriter ResultWriter) (err error) {
	searchTerms, err := setupSearchTerms(exportList)
	if err != nil {
		err = fmt.Errorf("setupSearchTerms() returned the following error: %w", err)
		return
	}

	volumeHandler, err := GetVolumeHandler(volumeLetter)
	if err != nil {
		err = fmt.Errorf("GetVolumeHandler() failed to get a handle to the volume %s: %w", volumeLetter, err)
		return
	}

	err = getFiles(volumeHandler, resultWriter, searchTerms)
	if err != nil {
		err = fmt.Errorf("getFiles() failed to get files: %w", err)
		return
	}

	return
}

func getFiles(volumeHandler VolumeHandler, resultWriter ResultWriter, listOfSearchKeywords listOfSearchTerms) (err error) {
	// Init a few things
	fileReaders := make(chan fileReader, 100)
	waitForInitialization := sync.WaitGroup{}
	waitForFileCopying := sync.WaitGroup{}

	waitForInitialization.Add(1)
	waitForFileCopying.Add(1)
	go func() {
		err = resultWriter.ResultWriter(&fileReaders, &waitForInitialization, &waitForFileCopying)
	}()
	waitForInitialization.Wait()
	if err != nil {
		err = fmt.Errorf("failed to setup resultWriter: %v", err)
		return
	}

	// parse the mft's mft record to get its dataruns
	mftRecord0, err := parseMFTRecord0(volumeHandler)
	if err != nil {
		err = fmt.Errorf("parseMFTRecord0() failed to parse mft record 0 from the volume %s: %w", volumeHandler.VolumeLetter, err)
		return
	}

	// Go back to the beginning of the mft record
	_, _ = syscall.Seek(volumeHandler.Handle, volumeHandler.Vbr.MftByteOffset, 0)

	// Open a raw reader on the MFT
	foundFile := foundFile{
		dataAttribute: mftRecord0.DataAttribute,
		fullPath:      "$mft",
	}
	mftReader := rawFileReader(volumeHandler, foundFile)

	// Do we need to stream a copy of the mft while we read it?
	areWeCopyingTheMFT := false
	for index, value := range listOfSearchKeywords {
		if value.fileNameString == "$mft" {
			areWeCopyingTheMFT = true

			// delete this from our search list
			listOfSearchKeywords[index] = listOfSearchKeywords[len(listOfSearchKeywords)-1]
			listOfSearchKeywords = listOfSearchKeywords[:len(listOfSearchKeywords)-1]

			break
		}
	}

	if areWeCopyingTheMFT == true {
		pipeReader, pipeWriter := io.Pipe()
		teeReader := io.TeeReader(mftReader, pipeWriter)
		fileReader := fileReader{
			fullPath: fmt.Sprintf("%s__$mft", volumeHandler.VolumeLetter),
			reader:   pipeReader,
		}
		fileReaders <- fileReader
		volumeHandler.mftReader = teeReader
		directoryTree, possibleMatches, err := findPossibleMatches(volumeHandler, listOfSearchKeywords)
		if err != nil {
			err = fmt.Errorf("findPossibleMatches() failed: %w", err)
			return
		}
		err = pipeWriter.Close()
		if err != nil {
			err = fmt.Errorf("failed to close writer pipe: %w", err)
			return
		}
	} else {
		volumeHandler.mftReader = mftReader
		directoryTree, possibleMatches, err := findPossibleMatches(volumeHandler, listOfSearchKeywords)
		if err != nil {
			err = fmt.Errorf("findPossibleMatches() failed: %w", err)
			return
		}
	}

	return
}

type possibleMatch struct {
	fileNameAttribute mft.FileNameAttribute
	dataAttribute     mft.DataAttribute
}

type possibleMatches []possibleMatch

func findPossibleMatches(volumeHandler VolumeHandler, listOfSearchKeywords listOfSearchTerms) (listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree, err error) {
	// Init memory
	unresolvedDirectorTree := make(mft.UnresolvedDirectoryTree)
	listOfPossibleMatches = make(possibleMatches, 0)

	for err != io.EOF {
		buffer := mft.RawMasterFileTableRecord(make([]byte, volumeHandler.Vbr.MftRecordSize))
		_, err = volumeHandler.mftReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				err = fmt.Errorf("findPossibleMatches() failed to read the mft: %w", err)
				return
			}
		}

		result, _ := buffer.IsThisAnMftRecord()
		if result == false {
			continue
		}

		result, err = buffer.IsThisADirectory()
		if result == true {
			unresolvedDirectory, _ := mft.ConvertRawMFTRecordToDirectory(buffer)
			unresolvedDirectorTree[unresolvedDirectory.RecordNumber] = unresolvedDirectory
		} else {
			// Parse what we need out of the entry for us to copy the file
			rawRecordHeader, _ := buffer.GetRawRecordHeader()
			recordHeader, _ := rawRecordHeader.Parse()
			rawAttributes, _ := buffer.GetRawAttributes(recordHeader)
			fileNameAttributes, _, dataAttribute, _ := rawAttributes.Parse(volumeHandler.Vbr.BytesPerCluster)
			for _, fileNameAttribute := range fileNameAttributes {
				if strings.Contains(fileNameAttribute.FileNamespace, "WIN32") == true || strings.Contains(fileNameAttribute.FileNamespace, "POSIX") {
					for _, value := range listOfSearchKeywords {
						if value.fileNameRegex != nil {
							if value.fileNameRegex.MatchString(strings.ToLower(fileNameAttribute.FileName)) == true {
								possibleMatch := possibleMatch{
									fileNameAttribute: fileNameAttribute,
									dataAttribute:     dataAttribute,
								}
								listOfPossibleMatches = append(listOfPossibleMatches, possibleMatch)
								break
							}
						} else {
							if value.fileNameString == strings.ToLower(fileNameAttribute.FileName) {

								break
							}
						}
					}
				}
				break
			}
		}
	}

	directoryTree, _ = unresolvedDirectorTree.Resolve(volumeHandler.VolumeLetter)
	return
}

func confirmFoundFiles(listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree) (foundFiles foundFiles, err error) {

	return
}

type foundFile struct {
	dataAttribute mft.DataAttribute
	fullPath      string
}

type foundFiles []foundFile
