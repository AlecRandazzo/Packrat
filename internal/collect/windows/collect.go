// Copyright (c) 2020 Alec Randazzo

package windows

import (
	"archive/zip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
)

// Collect will find and collect target files into a format depending on the resultWriter type
func Collect(exportList FileExportList, zipFile *zip.Writer) error {
	// volumeHandler as an arg is a dependency injection
	log.Debugf("Attempting to acquire the following files %+v", exportList)

	handler := NewVolumeHandler(strings.Trim(os.Getenv("SYSTEMDRIVE"), ":"))

	searchTerms, err := setupSearchTerms(exportList)
	if err != nil {
		err = fmt.Errorf("setupSearchTerms() returned the following error: %w", err)
		return err
	}

	err = handler.GetHandle()
	if err != nil {
		return err
	}
	defer handler.Handle().Close()

	var foundFiles foundFiles
	foundFiles, err = findFiles(handler, searchTerms, handler.vbr.BytesPerSector)
	if err != nil {
		return fmt.Errorf("findFiles() failed to find files: %w", err)
	}

	copyFiles(handler, foundFiles, zipFile)

	return nil
}

func findFiles(handler handler, listOfSearchKeywords searchTermsList, bytesPerSector uint) (foundFiles, error) {
	foundFiles := make(foundFiles, 0)

	// parse the mft's mft record to get its dataruns
	mftRecord0, err := parseMFTRecord0(handler)
	if err != nil {
		err = fmt.Errorf("parseMFTRecord0() failed to parse mft record 0 from the volume %s: %w", handler.VolumeLetter(), err)
		return nil, err
	}
	log.Debugf("Parsed the MFT's MFT record and got the following: %+v", mftRecord0)

	// Open a reader on the raw MFT
	foundFile := foundFile{
		dataRuns: mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns,
		fullPath: fmt.Sprintf(`%s:\$mft`, handler.VolumeLetter()),
	}
	handler.Handle().Seek(handler.Vbr().MftOffset, 0)
	mftReader := rawFileReader(handler, foundFile)
	log.Debug("Obtained a io.Reader to the MFT's dataruns.")

	handler.UpdateReader(mftReader)
	directoryTree := mft.DirectoryTree{}
	possibleMatches := possibleMatches{}
	possibleMatches, directoryTree, err = findPossibleMatches(handler, listOfSearchKeywords, bytesPerSector)
	if err != nil {
		err = fmt.Errorf("findPossibleMatches() failed: %w", err)
		return nil, err
	}

	foundFiles = confirmFoundFiles(listOfSearchKeywords, possibleMatches, directoryTree)
	if err != nil {
		err = fmt.Errorf("confirmFoundFiles() failed with error: %w", err)
		return nil, err
	}

	return foundFiles, nil
}

func copyFiles(handler handler, foundFiles foundFiles, zipFile *zip.Writer) {
	for _, file := range foundFiles {
		// try to get an io.reader via api first
		reader, err := apiFileReader(file)
		if err != nil {
			log.Debugf("Failed to get API handle, trying to get a raw io.Reader for '%s' with data runs: %+v", file.fullPath, file.dataRuns)
			// failed to get an API handle, trying to get an io.reader via raw method
			reader = rawFileReader(handler, file)
		} else {
			log.Debugf("Got an API io.Reader for '%s'.", file.fullPath)
		}
		normalizedFilePath := strings.ReplaceAll(file.fullPath, `\`, "_")
		normalizedFilePath = strings.ReplaceAll(normalizedFilePath, ":", "_")
		var writer io.Writer
		writer, err = zipFile.Create(normalizedFilePath)
		if err != nil {
			log.Errorf("failed to create file in zip: %v", err)
		}
		_, err = io.Copy(writer, reader)
		if err != nil {
			log.Errorf("failed to write %s", file.fullPath)
		}
	}
	return
}
