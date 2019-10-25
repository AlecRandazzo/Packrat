package windowscollector

import (
	"fmt"
	mft "github.com/AlecRandazzo/GoFor-MFT-Parser"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type possibleMatch struct {
	fileNameAttribute mft.FileNameAttribute
	dataAttribute     mft.DataAttribute
}

type possibleMatches []possibleMatch

func findPossibleMatches(volumeHandler VolumeHandler, listOfSearchKeywords listOfSearchTerms) (listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree, err error) {
	log.Debugf("Starting to scan the MFT's dataruns to create a tree of directories and to search for the for the following search terms: %+v", listOfSearchKeywords)

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
								if len(possibleMatch.dataAttribute.NonResidentDataAttribute.DataRuns) == 0 {
									log.Debugf("Found a possible file match for '%s' though it has nil data runs. Here's the hex for verification: %x", possibleMatch.fileNameAttribute.FileName, buffer)
								} else {
									log.Debugf("Found a possible file match: %+v. Here's the hex if you need to do debugging: %x", possibleMatch, buffer)
								}
								listOfPossibleMatches = append(listOfPossibleMatches, possibleMatch)
								break
							}
						} else {
							if value.fileNameString == strings.ToLower(fileNameAttribute.FileName) {
								possibleMatch := possibleMatch{
									fileNameAttribute: fileNameAttribute,
									dataAttribute:     dataAttribute,
								}
								if len(possibleMatch.dataAttribute.NonResidentDataAttribute.DataRuns) == 0 {
									log.Debugf("Found a possible file match for '%s' though it has nil data runs. Here's the hex for verification: %x", possibleMatch.fileNameAttribute.FileName, buffer)
								} else {
									log.Debugf("Found a possible file match: %+v. Here's the hex if you need to do debugging: %x", possibleMatch, buffer)
								}
								listOfPossibleMatches = append(listOfPossibleMatches, possibleMatch)
								break
							}
						}
					}
					break
				}
			}
		}
	}
	log.Debugf("Resolving %d directories we found to build their full paths.", len(unresolvedDirectorTree))
	directoryTree, _ = unresolvedDirectorTree.Resolve(volumeHandler.VolumeLetter)
	log.Debugf("Successfully resolved %d directories.", len(directoryTree))
	return
}

type foundFile struct {
	dataAttribute mft.DataAttribute
	fullPath      string
}

type foundFiles []foundFile

func confirmFoundFiles(listOfSearchKeywords listOfSearchTerms, listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree) (foundFilesList foundFiles, err error) {
	log.Debug("Determining what possible matches are true matches.")
	foundFilesList = make(foundFiles, 0)
	for _, possibleMatch := range listOfPossibleMatches {
		// First make sure that the parent directory is in the directory tree
		if _, ok := directoryTree[possibleMatch.fileNameAttribute.ParentDirRecordNumber]; ok {
			// check against all the list of possible full paths
			possibleMatchFullPath := fmt.Sprintf("%s\\%s", strings.ToLower(directoryTree[possibleMatch.fileNameAttribute.ParentDirRecordNumber]), strings.ToLower(possibleMatch.fileNameAttribute.FileName))
			numberOfSearchTerms := len(listOfSearchKeywords)
			counter := 0
			for _, searchTerms := range listOfSearchKeywords {
				if searchTerms.fullPathRegex != nil {
					if searchTerms.fullPathRegex.MatchString(possibleMatchFullPath) == true {
						foundFile := foundFile{
							dataAttribute: possibleMatch.dataAttribute,
							fullPath:      possibleMatchFullPath,
						}
						log.Debugf("Found a true match: %+v", foundFile)
						foundFilesList = append(foundFilesList, foundFile)
						break
					}
				} else {
					if searchTerms.fullPathString == possibleMatchFullPath {
						foundFile := foundFile{
							dataAttribute: possibleMatch.dataAttribute,
							fullPath:      possibleMatchFullPath,
						}
						log.Debugf("Found a true match: %+v", foundFile)
						foundFilesList = append(foundFilesList, foundFile)
						break
					}
				}
				counter++
				if counter == numberOfSearchTerms {
					log.Debugf("The file %s did end up being a true positive", possibleMatchFullPath)
				}
			}
		} else {
			// continue if parent directory is not in the directory tree map
			continue
		}
	}
	return
}
