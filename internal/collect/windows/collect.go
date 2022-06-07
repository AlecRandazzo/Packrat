// Copyright (c) 2022 Alec Randazzo

package windows

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecRandazzo/Packrat/internal/collect"
	"github.com/AlecRandazzo/Packrat/pkg/windows/mft"
	"github.com/AlecRandazzo/Packrat/pkg/windows/volume"
	log "github.com/sirupsen/logrus"
)

// Collect will find and collect target files into a format depending on the resultWriter type
func Collect(exportList FileExportList, writer collect.Writer) error {
	// volumeHandler as an arg is a dependency injection
	log.Debugf("Attempting to acquire the following files %+v", exportList)

	handler, err := volume.NewHandler(strings.Trim(os.Getenv("SYSTEMDRIVE"), ":"))
	if err != nil {
		return fmt.Errorf("failed to get volume handle: %w", err)
	}

	var searchTerms searchTermsList
	searchTerms, err = setupSearchTerms(exportList)
	if err != nil {
		err = fmt.Errorf("setupSearchTerms() returned the following error: %w", err)
		return err
	}

	defer handler.Handle().Close()

	var foundFiles foundFiles
	foundFiles, err = findFiles(handler, searchTerms)
	if err != nil {
		return fmt.Errorf("findFiles() failed to find files: %w", err)
	}

	copyFiles(handler, foundFiles, writer)

	return nil
}

// findFiles searches for files that we want to collect
func findFiles(handler volume.Handler, listOfSearchKeywords searchTermsList) (foundFiles, error) {
	var foundFiles foundFiles
	// parser the mft's mft record to get its dataruns
	mftRecord0, err := parseMFTRecord0(handler)
	if err != nil {
		err = fmt.Errorf("parseMFTRecord0() failed to parser mft record 0 from the volume %s: %w", handler.Letter(), err)
		return nil, err
	}
	log.Debugf("Parsed the MFT's MFT record and got the following: %+v", mftRecord0)

	// Open a reader on the raw MFT
	foundFile := foundFile{
		dataRuns: mftRecord0.DataAttribute.NonResidentDataAttribute.DataRuns,
		fullPath: fmt.Sprintf(`%s:\$mft`, handler.Letter()),
	}
	_, _ = handler.Handle().Seek(handler.Vbr().MftOffset, 0)
	mftReader := newRawFileReader(handler, foundFile)
	log.Debug("Obtained a io.Reader to the MFT's dataruns.")

	handler.UpdateReader(mftReader)
	var directoryTree mft.DirectoryTree
	var possibleMatches possibleMatches
	possibleMatches, directoryTree, err = findPossibleMatches(handler, listOfSearchKeywords)
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

// copyFiles copies the files we want to collect to whatever we set as our collect.Writer.
func copyFiles(handler volume.Handler, foundFiles foundFiles, writer collect.Writer) {
	for _, file := range foundFiles {
		// try to get an io.reader via api first
		reader, err := newApiFileReader(file)
		if err != nil {
			log.Debugf("Failed to get API handle, trying to get a raw io.Reader for '%s' with data runs: %+v", file.fullPath, file.dataRuns)
			// failed to get an API handle, trying to get an io.reader via raw method
			reader = newRawFileReader(handler, file)
		} else {
			log.Debugf("Got an API io.Reader for '%s'.", file.fullPath)
		}

		normalizedFilePath := normalizeFilePath(file.fullPath)
		err = writer.Write(reader, normalizedFilePath)
		if err != nil {
			log.Errorf("failed to write %s to zip: %s", file.fullPath, err.Error())
		}
	}
}

// normalizeFilePath replaces "\" and ":" with an "_"
func normalizeFilePath(filepath string) string {
	normalizeFilePath := strings.ReplaceAll(filepath, `\`, "_")
	normalizeFilePath = strings.ReplaceAll(normalizeFilePath, ":", "_")
	return normalizeFilePath
}
