/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package windowscollector

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// File that you want to export.
type FileToExport struct {
	FullPath        string
	IsFullPathRegex bool
	FileName        string
	IsFileNameRegex bool
}

// Slice of files that you want to export.
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
