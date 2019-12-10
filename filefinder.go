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
	mft "github.com/Go-Forensics/MFT-Parser"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
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

func checkForPossibleMatch(listOfSearchKeywords listOfSearchTerms, fileNameAttributes mft.FileNameAttributes) (result bool, fileNameAttribute mft.FileNameAttribute, err error) {
	// Sanity Checking
	if len(listOfSearchKeywords) == 0 {
		err = errors.New("checkForPossibleMatch() received an empty listOfSearchTerms")
		return
	}
	if len(fileNameAttributes) == 0 {
		err = errors.New("checkForPossibleMatch() received an empty fileNameAttributes")
		return
	}

	for _, attribute := range fileNameAttributes {
		if strings.Contains(attribute.FileNamespace, "WIN32") == true || strings.Contains(attribute.FileNamespace, "POSIX") {
			for _, value := range listOfSearchKeywords {
				if value.fileNameRegex != nil {
					if value.fileNameRegex.MatchString(strings.ToLower(attribute.FileName)) == true {
						result = true
						fileNameAttribute = attribute
						return
					}
				} else {
					if value.fileNameString == strings.ToLower(attribute.FileName) {
						result = true
						fileNameAttribute = attribute
						return
					}
				}
			}
		}
	}

	result = false
	return
}

func findPossibleMatches(volumeHandler *VolumeHandler, listOfSearchKeywords listOfSearchTerms) (listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree, err error) {
	log.Debugf("Starting to scan the MFT's dataruns to create a tree of directories and to search for the for the following search terms: %+v", listOfSearchKeywords)

	// Init memory
	unresolvedDirectorTree := make(mft.UnresolvedDirectoryTree)
	listOfPossibleMatches = make(possibleMatches, 0)
	recordOffsetTracker := make(mftRecordVolumeOffsetTracker)
	listOfMftRecordWithNonResidentAttributes := make(listOfMftRecordWithNonResidentAttributes, 0)

	for err != io.EOF {
		buffer := mft.RawMasterFileTableRecord(make([]byte, volumeHandler.Vbr.MftRecordSize))
		_, err = volumeHandler.mftReader.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}

		result, _ := buffer.IsThisAnMftRecord()
		if result == false {
			continue
		}

		result, err = buffer.IsThisADirectory()
		if result == true {
			unresolvedDirectory, _ := mft.ConvertRawMFTRecordToDirectory(buffer)
			unresolvedDirectorTree[unresolvedDirectory.RecordNumber] = unresolvedDirectory
			recordOffsetTracker[unresolvedDirectory.RecordNumber] = volumeHandler.lastReadVolumeOffset
		} else {
			// Parse what we need out of the entry for us to copy the file
			rawRecordHeader, _ := buffer.GetRawRecordHeader()
			recordHeader, _ := rawRecordHeader.Parse()
			recordOffsetTracker[recordHeader.RecordNumber] = volumeHandler.lastReadVolumeOffset
			rawAttributes, _ := buffer.GetRawAttributes(recordHeader)
			fileNameAttributes, _, dataAttribute, attributeListAttributes, _ := rawAttributes.Parse(volumeHandler.Vbr.BytesPerCluster)
			result, fileNameAttribute, err := checkForPossibleMatch(listOfSearchKeywords, fileNameAttributes)
			if err != nil || result == false {
				continue
			}

			if attributeListAttributes == nil {
				log.Debugf("Found a possible match. File name is '%s' and its MFT offset is %d. Here is the MFT record hex: %x", fileNameAttribute.FileName, volumeHandler.lastReadVolumeOffset, []byte(buffer))
				aPossibleMatch := possibleMatch{
					fileNameAttribute: fileNameAttribute,
					dataRuns:          dataAttribute.NonResidentDataAttribute.DataRuns,
				}
				listOfPossibleMatches = append(listOfPossibleMatches, aPossibleMatch)
				continue
			} else {
				log.Debugf("Found a possible match which has an attribute list. File name is '%s' and its MFT offset is %d. Here is the attribute list: %+v Here is the MFT record hex: %x", fileNameAttribute.FileName, volumeHandler.lastReadVolumeOffset, attributeListAttributes, buffer)
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
		newVolumeHandle, _ := volumeHandler.GetHandle(volumeHandler.VolumeLetter)
		for _, record := range listOfMftRecordWithNonResidentAttributes {
			attributeCounter := 0
			sizeOfAttributeListAttributes := len(record.attributeListAttributes)
			dataRuns := make(mft.DataRuns)
			for attributeCounter < sizeOfAttributeListAttributes {
				switch record.attributeListAttributes[attributeCounter].Type {
				case 0x80:
					nonResidentRecordNumber := record.attributeListAttributes[attributeCounter].MFTReferenceRecordNumber
					absoluteVolumeOffset := recordOffsetTracker[nonResidentRecordNumber]
					_, _ = newVolumeHandle.Seek(absoluteVolumeOffset, 0)
					buffer := mft.RawMasterFileTableRecord(make([]byte, volumeHandler.Vbr.BytesPerCluster))
					_, _ = newVolumeHandle.Read(buffer)
					mftRecord, _ := buffer.Parse(volumeHandler.Vbr.BytesPerCluster)
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
				dataRuns:          dataRuns,
			}
			log.Debugf("Pieced together a series of non resident data attributes and got the following: %+v", aPossibleMatch)
			listOfPossibleMatches = append(listOfPossibleMatches, aPossibleMatch)
		}
	}

	log.Debugf("Resolving %d directories we found to build their full paths.", len(unresolvedDirectorTree))
	directoryTree, _ = unresolvedDirectorTree.Resolve(volumeHandler.VolumeLetter)
	log.Debugf("Successfully resolved %d directories.", len(directoryTree))
	return
}

type foundFile struct {
	dataRuns mft.DataRuns
	fullPath string
	fileSize int64
}

type foundFiles []foundFile

func confirmFoundFiles(listOfSearchKeywords listOfSearchTerms, listOfPossibleMatches possibleMatches, directoryTree mft.DirectoryTree) (foundFilesList foundFiles) {
	log.Debug("Determining what possible matches are true matches.")
	foundFilesList = make(foundFiles, 0)
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
							fileSize: int64(possibleMatch.fileNameAttribute.PhysicalFileSize),
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
	return
}
