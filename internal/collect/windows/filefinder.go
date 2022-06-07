// Copyright (c) 2022 Alec Randazzo

package windows

import (
	"errors"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/windows/mft"
)

type possibleMatch struct {
	fileNameAttribute mft.FileNameAttribute
	dataRuns          mft.DataRuns
}

type possibleMatches []possibleMatch

type mftRecordVolumeOffsetTracker map[uint32]int64

type mftRecordWithNonResidentAttributes struct {
	fnAttribute             mft.FileNameAttribute
	dataAttribute           mft.DataAttribute
	attributeListAttributes mft.AttributeListAttributes
}

type listOfMftRecordWithNonResidentAttributes []mftRecordWithNonResidentAttributes

func checkForPossibleMatch(listOfSearchKeywords searchTermsList, fileNameAttributes mft.FileNameAttributes) (mft.FileNameAttribute, error) {
	// Sanity Checking
	if len(listOfSearchKeywords) == 0 {
		return mft.FileNameAttribute{}, errors.New("checkForPossibleMatch() received an empty searchTermsList")
	}
	if len(fileNameAttributes) == 0 {
		return mft.FileNameAttribute{}, errors.New("checkForPossibleMatch() received an empty fileNameAttributes")
	}

	for _, attribute := range fileNameAttributes {
		if strings.Contains(string(attribute.FileNamespace), "WIN32") ||
			strings.Contains(string(attribute.FileNamespace), "WIN32 & DOS") ||
			strings.Contains(string(attribute.FileNamespace), "POSIX") {
			for _, value := range listOfSearchKeywords {
				if value.fileNameRegex != nil {
					if value.fileNameRegex.MatchString(strings.ToLower(attribute.FileName)) {
						return attribute, nil
					}
				} else {
					if value.fileNameString == strings.ToLower(attribute.FileName) {
						return attribute, nil
					}
				}
			}
		}
	}

	return mft.FileNameAttribute{}, errors.New("no match")
}

func findPossibleMatches(handler handler, listOfSearchKeywords searchTermsList, bytesPerSector uint) (possibleMatches, mft.DirectoryTree, error) {
	log.Debugf("Starting to scan the MFT's dataruns to create a tree of directories and to search for the for the following search terms: %+v", listOfSearchKeywords)

	// Init memory
	unresolvedDirectorTree := make(mft.UnresolvedDirectoryTree)
	listOfPossibleMatches := make(possibleMatches, 0)
	recordOffsetTracker := make(mftRecordVolumeOffsetTracker)
	listOfMftRecordWithNonResidentAttributes := make(listOfMftRecordWithNonResidentAttributes, 0)
	var err error

	for err != io.EOF {
		buffer := make([]byte, handler.Vbr().MftRecordSize)
		_, err = handler.Reader().Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}

		err = mft.ValidateMftRecordBytes(buffer)
		if err != nil {
			continue
		}

		err = mft.ValidateDirectory(buffer)
		if err == nil {
			unresolvedDirectory, _ := mft.ConvertRawMFTRecordToDirectory(buffer, bytesPerSector)
			unresolvedDirectorTree[unresolvedDirectory.RecordNumber] = unresolvedDirectory
			recordOffsetTracker[unresolvedDirectory.RecordNumber] = handler.LastOffset()
		} else {
			// Parse what we need out of the entry for us to copy the file
			recordHeader, _ := mft.GetRecordHeaders(buffer)
			recordOffsetTracker[recordHeader.RecordNumber] = handler.LastOffset()
			rawAttributes, _ := mft.GetRawAttributes(buffer, recordHeader)
			fileNameAttributes, _, dataAttribute, attributeListAttributes, _ := mft.GetAttributes(rawAttributes, handler.Vbr().BytesPerCluster)
			var fileNameAttribute mft.FileNameAttribute
			fileNameAttribute, err = checkForPossibleMatch(listOfSearchKeywords, fileNameAttributes)
			if err != nil {
				continue
			}

			if attributeListAttributes == nil {
				log.Debugf("Found a possible match. File name is '%s' and its MFT offset is %d. Here is the MFT record hex: %x", fileNameAttribute.FileName, handler.LastOffset(), buffer)
				aPossibleMatch := possibleMatch{
					fileNameAttribute: fileNameAttribute,
					dataRuns:          dataAttribute.NonResidentDataAttribute.DataRuns,
				}
				listOfPossibleMatches = append(listOfPossibleMatches, aPossibleMatch)
				continue
			} else {
				log.Debugf("Found a possible match which has an attribute list. File name is '%s' and its MFT offset is %d. Here is the attribute list: %+v Here is the MFT record hex: %x", fileNameAttribute.FileName, handler.LastOffset(), attributeListAttributes, buffer)
				trackThisForLater := mftRecordWithNonResidentAttributes{
					fnAttribute:             fileNameAttribute,
					dataAttribute:           dataAttribute,
					attributeListAttributes: attributeListAttributes,
				}
				listOfMftRecordWithNonResidentAttributes = append(listOfMftRecordWithNonResidentAttributes, trackThisForLater)
				continue
			}
		}
	}

	// Resolve the possible matches that had attribute lists
	if len(listOfMftRecordWithNonResidentAttributes) != 0 {
		newVolumeHandle := NewVolumeHandler(handler.VolumeLetter())
		_ = newVolumeHandle.GetHandle()
		for _, record := range listOfMftRecordWithNonResidentAttributes {
			attributeCounter := 0
			sizeOfAttributeListAttributes := len(record.attributeListAttributes)
			dataRuns := make(mft.DataRuns)
			for attributeCounter < sizeOfAttributeListAttributes {
				switch record.attributeListAttributes[attributeCounter].Type {
				case 0x80:
					nonResidentRecordNumber := record.attributeListAttributes[attributeCounter].MFTReferenceRecordNumber
					absoluteVolumeOffset := recordOffsetTracker[nonResidentRecordNumber]
					_, _ = newVolumeHandle.Handle().Seek(absoluteVolumeOffset, 0)
					buffer := make([]byte, handler.Vbr().BytesPerCluster)
					_, _ = newVolumeHandle.Handle().Read(buffer)
					mftRecord, _ := mft.ParseRecord(buffer, handler.Vbr().BytesPerSector, handler.Vbr().BytesPerCluster)
					log.Debugf("Went to absolute offset %d to get a non resident data attribute with record number %d. Parsed the record for the values %+v. Raw hex: %x", absoluteVolumeOffset, nonResidentRecordNumber, mftRecord, buffer)
					tempDataRunCounter := 0
					numberOfDataRuns := len(mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns)
					for tempDataRunCounter < numberOfDataRuns {
						index := len(dataRuns)
						dataRuns[index] = mftRecord.DataAttribute.NonResidentDataAttribute.DataRuns[tempDataRunCounter]
						tempDataRunCounter++
					}
					attributeCounter++
				default:
					attributeCounter++
				}
			}
			aPossibleMatch := possibleMatch{
				fileNameAttribute: record.fnAttribute,
				dataRuns:          record.dataAttribute.NonResidentDataAttribute.DataRuns,
			}
			log.Debugf("Pieced together a series of non resident data attributes and got the following: %+v", aPossibleMatch)
			listOfPossibleMatches = append(listOfPossibleMatches, aPossibleMatch)
		}
	}

	log.Debugf("Resolving %d directories we found to build their full paths.", len(unresolvedDirectorTree))
	directoryTree, _ := unresolvedDirectorTree.Resolve(handler.VolumeLetter())
	log.Debugf("Successfully resolved %d directories.", len(directoryTree))

	return listOfPossibleMatches, directoryTree, nil
}

type foundFile struct {
	dataRuns mft.DataRuns
	fullPath string
	size     int64
}

type foundFiles []foundFile

func confirmFoundFiles(listOfSearchKeywords searchTermsList, listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree) foundFiles {
	log.Debug("Determining what possible matches are true matches.")
	foundFilesList := make(foundFiles, 0)
	for _, possibleMatch := range listOfPossibleMatches {
		// First make sure that the parent directory is in the directory tree
		if _, ok := directoryTree[possibleMatch.fileNameAttribute.ParentDirRecordNumber]; ok {
			// check against all the list of possible full paths
			possibleMatchFullPath := fmt.Sprintf(`%s\%s`, strings.ToLower(directoryTree[possibleMatch.fileNameAttribute.ParentDirRecordNumber]), strings.ToLower(possibleMatch.fileNameAttribute.FileName))
			numberOfSearchTerms := len(listOfSearchKeywords)
			counter := 0
			for _, searchTerms := range listOfSearchKeywords {
				if searchTerms.fullPathRegex != nil {
					if searchTerms.fullPathRegex.MatchString(possibleMatchFullPath) == true {
						foundFile := foundFile{
							dataRuns: possibleMatch.dataRuns,
							fullPath: possibleMatchFullPath,
							size:     int64(possibleMatch.fileNameAttribute.PhysicalFileSize),
						}
						log.Debugf("Found a true match: %+v", foundFile)
						foundFilesList = append(foundFilesList, foundFile)
						break
					}
				} else {
					if searchTerms.fullPathString == possibleMatchFullPath {
						foundFile := foundFile{
							dataRuns: possibleMatch.dataRuns,
							fullPath: possibleMatchFullPath,
						}
						log.Debugf("Found a true match: %+v", foundFile)
						foundFilesList = append(foundFilesList, foundFile)
						break
					}
				}
				counter++
				if counter == numberOfSearchTerms {
					log.Debugf("The file %s did not end up being a true positive", possibleMatchFullPath)
				}
			}
		} else {
			// continue if parent directory is not in the directory tree map
			continue
		}
	}
	return foundFilesList
}
