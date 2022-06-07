// Copyright (c) 2022 Alec Randazzo

package windows

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// FileExport is the file that you want to export.
type FileExport struct {
	FullPath      string
	FullPathRegex bool
	FileName      string
	FileNameRegex bool
}

// FileExportList is a slice of files that you want to export.
type FileExportList []FileExport

type searchTerms struct {
	fullPathString string
	fullPathRegex  *regexp.Regexp
	fileNameString string
	fileNameRegex  *regexp.Regexp
}

type searchTermsList []searchTerms

// setupSearchTerms sets up search terms into a struct that is used by other functions
func setupSearchTerms(exportList FileExportList) (searchTermsList, error) {
	listOfSearchKeywords := make(searchTermsList, 0)
	for _, value := range exportList {
		// Sanity checking inputs
		if value.FileName == "" {
			return nil, errors.New("received empty filename string")
		} else if value.FullPath == "" {
			return nil, errors.New("received empty filepath string")
		}

		// Normalize everything
		value.FullPath = strings.ToLower(value.FullPath)
		value.FileName = strings.ToLower(value.FileName)

		if !value.FullPathRegex && strings.HasSuffix(value.FullPath, `\`) {
			return nil, fmt.Errorf("file path '%s' has a trailing '\\'", value.FullPath)
		} else if !value.FullPathRegex && strings.HasSuffix(value.FullPath, `\`) {
			return nil, fmt.Errorf("file path '%s' has missing a trailing '\\\\'", value.FullPath)
		}

		searchKeywords := searchTerms{}
		switch value.FullPathRegex {
		case false:
			searchKeywords.fullPathString = value.FullPath
			searchKeywords.fullPathRegex = nil
		case true:
			searchKeywords.fullPathString = ""
			var err error
			searchKeywords.fullPathRegex, err = regexp.Compile(value.FullPath)
			if err != nil {
				return nil, fmt.Errorf("failed to compile regex for %s: %w", value.FullPath, err)
			}
		}

		switch value.FileNameRegex {
		case false:
			searchKeywords.fileNameString = value.FileName
			searchKeywords.fileNameRegex = nil
		case true:
			searchKeywords.fileNameString = ""
			searchKeywords.fileNameRegex = regexp.MustCompile(value.FileName)
		}

		listOfSearchKeywords = append(listOfSearchKeywords, searchKeywords)
	}

	return listOfSearchKeywords, nil
}
