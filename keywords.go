// Copyright (c) 2020 Alec Randazzo

package windowscollector

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// FileToExport is the file that you want to export.
type FileToExport struct {
	FullPath        string `yaml:"FullPath"`
	IsFullPathRegex bool   `yaml:"IsFullPathRegex"`
	FileName        string `yaml:"FileName"`
	IsFileNameRegex bool   `yaml:"IsFileNameRegex"`
}

// ListOfFilesToExport is a slice of files that you want to export.
type ListOfFilesToExport []FileToExport

type searchTerms struct {
	fullPathString string
	fullPathRegex  *regexp.Regexp
	fileNameString string
	fileNameRegex  *regexp.Regexp
}

type listOfSearchTerms []searchTerms

func setupSearchTerms(exportList ListOfFilesToExport) (listOfSearchKeywords listOfSearchTerms, err error) {
	for _, value := range exportList {
		// Sanity checking inputs
		if value.FileName == "" {
			err = errors.New("received empty filename string")
			return
		} else if value.FullPath == "" {
			err = errors.New("received empty filepath string")
			return
		}

		// Normalize everything
		value.FullPath = strings.ToLower(value.FullPath)
		value.FileName = strings.ToLower(value.FileName)

		if value.IsFullPathRegex == false && strings.HasSuffix(value.FullPath, `\`) == true {
			err = fmt.Errorf("file path '%s' has a trailing '\\'", value.FullPath)
			return
		} else if value.IsFullPathRegex == true && strings.HasSuffix(value.FullPath, `\`) == true {
			err = fmt.Errorf("file path '%s' has missing a trailing '\\\\'", value.FullPath)
			return
		}

		searchKeywords := searchTerms{}
		switch value.IsFullPathRegex {
		case false:
			searchKeywords.fullPathString = value.FullPath
			searchKeywords.fullPathRegex = nil
		case true:
			searchKeywords.fullPathString = ""
			searchKeywords.fullPathRegex = regexp.MustCompile(value.FullPath)
		}

		switch value.IsFileNameRegex {
		case false:
			searchKeywords.fileNameString = value.FileName
			searchKeywords.fileNameRegex = nil
		case true:
			searchKeywords.fileNameString = ""
			searchKeywords.fileNameRegex = regexp.MustCompile(value.FileName)
		}

		listOfSearchKeywords = append(listOfSearchKeywords, searchKeywords)
	}

	return
}
