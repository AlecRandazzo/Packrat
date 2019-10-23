package windowscollector

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// File that you want to export.
type FileToExport struct {
	FilePath        string
	IsFilePathRegex bool
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
	filePathString string
	filePathRegex  *regexp.Regexp
}

type listOfSearchTerms []searchTerms

func setupSearchTerms(exportList ListOfFilesToExport) (listOfSearchKeywords listOfSearchTerms, err error) {
	for _, value := range exportList {
		// Sanity checking inputs
		if value.FileName == "" {
			err = errors.New("received empty filename string")
			return
		} else if value.FilePath == "" {
			err = errors.New("received empty filepath string")
			return
		}

		// Normalize everything
		re := regexp.MustCompile("^.:")
		value.FilePath = strings.ToLower(re.ReplaceAllString(value.FilePath, ":"))
		value.FileName = strings.ToLower(value.FileName)

		if value.IsFilePathRegex == false && strings.HasSuffix(value.FilePath, "\\") == true {
			err = fmt.Errorf("file path '%s' has a trailing '\\'", value.FilePath)
			return
		} else if value.IsFilePathRegex == true && strings.HasSuffix(value.FilePath, "\\\\") == true {
			err = fmt.Errorf("file path '%s' has missing a trailing '\\\\'", value.FilePath)
			return
		}

		searchKeywords := searchTerms{}
		switch value.IsFilePathRegex {
		case false:
			searchKeywords.filePathString = value.FilePath
			searchKeywords.filePathRegex = nil
		case true:
			searchKeywords.filePathString = ""
			searchKeywords.filePathRegex = regexp.MustCompile(value.FilePath)
		}

		switch value.IsFileNameRegex {
		case false:
			searchKeywords.fileNameString = value.FileName
			searchKeywords.fileNameRegex = nil
		case true:
			searchKeywords.fileNameString = ""
			searchKeywords.fileNameRegex = regexp.MustCompile(value.FileName)
		}

		if value.IsFileNameRegex == false && value.IsFilePathRegex == false {
			searchKeywords.fullPathString = value.FilePath + value.FileName
			searchKeywords.fullPathRegex = nil
		} else {
			searchKeywords.fullPathString = ""
			convertThisToRegex := value.FilePath + value.FileName
			searchKeywords.fullPathRegex = regexp.MustCompile(convertThisToRegex)
		}

		listOfSearchKeywords = append(listOfSearchKeywords, searchKeywords)
	}

	return
}
