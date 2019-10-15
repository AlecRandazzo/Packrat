package windowscollector

import (
	"fmt"
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/windows"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

func Collect(volumeLetter string, exportList ExportList, resultWriter ResultWriter) (err error) {
	volumeHandler, err := GetVolumeHandler(volumeLetter)
	if err != nil {
		err = fmt.Errorf("failed to get a handle to the volume %v: %v", volumeLetter, err)
		return
	}

	mftRecord0, err := parseMFTRecord0(volumeHandler)
	if err != nil {
		err = fmt.Errorf("failed to parse mft record 0 from the volume %v: %v", volumeLetter, err)
		return
	}

	directoryTree, err := buildDirectoryTree(volumeHandler, mftRecord0)
	if err != nil {
		err = fmt.Errorf("failed to build a directory tree for volume %v: %v", volumeLetter, err)
		return
	}

	fileEqualSearchList, fileRegexSearchList := buildFileExportLists(exportList)

	foundFiles := findFiles(volumeHandler, mftRecord0, directoryTree, fileEqualSearchList, fileRegexSearchList)

	err = copyFiles(volumeHandler, foundFiles, resultWriter)

	return
}

func copyFiles(volumeHandler VolumeHandler, foundFiles foundFiles, resultWriter ResultWriter) (err error) {
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

	for _, file := range foundFiles {
		var reader io.Reader

		// Defang some special characters int he full path string
		fullPath := strings.ReplaceAll(file.fullPath, "\\", "_")
		fullPath = strings.ReplaceAll(fullPath, ":", "_")

		// Try to copy the file via API first. If it fails, try to copy the file via raw.
		reader, err = apiCopy(file)
		if err != nil {
			reader = rawCopy(volumeHandler, file)
		}
		fileReaders <- fileReader{
			fullPath: fullPath,
			reader:   reader,
		}
	}
	waitForFileCopying.Wait()
	return
}

func apiCopy(file foundFile) (reader io.Reader, err error) {
	reader, err = os.Open(file.fullPath)
	return
}

func rawCopy(handler VolumeHandler, file foundFile) (reader io.Reader) {
	reader = &DataRunsReader{
		VolumeHandler:          handler,
		DataRuns:               file.mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns,
		fileName:               file.fullPath,
		dataRunTracker:         0,
		bytesLeftToReadTracker: 0,
		initialized:            false,
	}
	return
}

type foundFile struct {
	mftRecord mft.MasterFileTableRecord
	fullPath  string
}

type foundFiles []foundFile

func findFiles(volumeHandler VolumeHandler, mftRecord0 mft.MasterFileTableRecord, directoryTree mft.DirectoryTree, fileEqualSearchList fileEqualSearchList, fileRegexSearchList fileRegexSearchList) (foundFiles foundFiles) {
	// Init memory
	foundFiles = make([]foundFile, 0)

	// Search mft record 0's data runs for the files we want to collect
	for _, dataRun := range mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns {
		// Seek to the start of the datarun
		_, _ = syscall.Seek(volumeHandler.Handle, dataRun.AbsoluteOffset, 0)

		// Keep track of how large the data run is. We will be counting down this number so we don't overshoot.
		bytesLeft := dataRun.Length

		// This for loop will run while there is datarun left to read.
		for bytesLeft > 0 {
			// Load bytes into a buffer
			buffer := make([]byte, volumeHandler.Vbr.MftRecordSize)
			_, _ = syscall.Read(volumeHandler.Handle, buffer)

			// Parse the buffered bytes
			mftRecord, _ := mft.RawMasterFileTableRecord(buffer).Parse(volumeHandler.Vbr.BytesPerCluster)
			if len(mftRecord.FileNameAttributes) == 0 {
				bytesLeft -= volumeHandler.Vbr.MftRecordSize
				continue
			}

			for _, attribute := range mftRecord.FileNameAttributes {
				// Find the filename attribute that will have the non DOS name for the file
				if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") == true {
					recordFullPath := strings.ToLower(directoryTree[attribute.ParentDirRecordNumber] + "\\" + attribute.FileName)

					// Cross check our string search list. If found, track the file for collection later
					for _, file := range fileEqualSearchList {
						if file == recordFullPath {
							log.Debugf("Found the MFT record for the file '%s', tracking relevant metadata for copying later.", recordFullPath)
							foundFile := foundFile{
								mftRecord: mftRecord,
								fullPath:  recordFullPath,
							}
							foundFiles = append(foundFiles, foundFile)
							break
						}
					}

					// Cross check our regex search list. If found, track the file for collection later
					for _, file := range fileRegexSearchList {
						if file.Match([]byte(recordFullPath)) == true {
							log.Debugf("Found the MFT record for the file '%s', tracking relevant metadata for copying later.", recordFullPath)
							foundFile := foundFile{
								mftRecord: mftRecord,
								fullPath:  recordFullPath,
							}
							foundFiles = append(foundFiles, foundFile)
							break
						}
					}
					break
				}
			}

			// Reduce our byte tracker by the number of bytes read
			bytesLeft -= volumeHandler.Vbr.MftRecordSize
		}
	}

	return
}

// File that you want to export.
type FileToExport struct {
	FullPath string
	Type     string
}

// Slice of files that you want to export.
type ExportList []FileToExport
type fileEqualSearchList []string
type fileRegexSearchList []*regexp.Regexp

func buildFileExportLists(exportList ExportList) (fileEqualList fileEqualSearchList, fileRegexList fileRegexSearchList) {
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
