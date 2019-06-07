/*
 * Copyright (c) 2019 Alec Randazzo
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 */

package gofor

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

type ResidentDataAttributes struct {
	ResidentData []byte
}

type NonResidentDataAttributes struct {
	StartingVCN   int
	EndingVCN     int
	OffsetDataRun int8
	AllocatedSize uint64
	RealSize      uint64
	DataRuns      map[int]DataRun
}

type RawDataRun struct {
	NumberOrder      int
	ClusterOffset    int64
	NumberOfClusters int64
}

type DataRun struct {
	AbsoluteOffset int64
	Length         int64
}

type DataAttributes struct {
	TotalSize                 uint8
	FlagResident              bool
	ResidentDataAttributes    ResidentDataAttributes
	NonResidentDataAttributes NonResidentDataAttributes
}

func (mftRecord *MasterFileTableRecord) getDataAttribute() (err error) {
	const codeData = 0x80
	const offsetResidentFlag = 0x08

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("Failed to parse data attribute")
		}
	}()

	for _, attribute := range mftRecord.AttributeInfo {
		if attribute.AttributeType == codeData {
			if len(attribute.AttributeBytes) < 0x18 {
				return
			}

			//TODO: handle resident data
			if attribute.AttributeBytes[offsetResidentFlag] == 0x00 {
				mftRecord.DataAttributes.FlagResident = true
				mftRecord.DataAttributes.ResidentDataAttributes, err = getResidentDataAttribute(attribute.AttributeBytes)
				if err != nil {
					err = errors.Wrap(err, "failed to parse resident data attribute")
					return
				}
				return
			} else {
				mftRecord.DataAttributes.FlagResident = false
				mftRecord.DataAttributes.NonResidentDataAttributes, err = getNonResidentDataAttribute(attribute.AttributeBytes, mftRecord.bytesPerCluster)
				if err != nil {
					err = errors.Wrap(err, "failed to parse non resident data attribute")
					return
				}
			}
			break
		}
	}
	return
}

func getResidentDataAttribute(attributeBytes []byte) (residentDataAttributes ResidentDataAttributes, err error) {
	const offsetResidentData = 0x18

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovery %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(attributeBytes))
		}
	}()

	if len(attributeBytes) < 0x18 {
		return
	}

	residentDataAttributes.ResidentData = attributeBytes[offsetResidentData:]

	return
}

func getNonResidentDataAttribute(attributeBytes []byte, bytesPerCluster int64) (nonResidentDataAttributes NonResidentDataAttributes, err error) {
	const offsetDataRunOffset = 0x20

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovery %s, hex dump: %s", fmt.Sprint(r), hex.EncodeToString(attributeBytes))
		}
	}()

	if len(attributeBytes) <= 0x20 {
		return
	}

	// Identify offset of the data runs in the data attribute
	dataRunOffset := attributeBytes[offsetDataRunOffset]

	if len(attributeBytes) < int(dataRunOffset) {
		return
	}

	// Pull out the data run bytes
	dataRunsBytes := attributeBytes[dataRunOffset:]

	// Send the bytes to be parsed
	nonResidentDataAttributes.DataRuns = getDataRuns(dataRunsBytes, bytesPerCluster)
	if nonResidentDataAttributes.DataRuns == nil {
		err = errors.Wrap(err, "failed to identify data runs")
		return
	}

	return
}

func resolveDataRuns(rawDataRuns map[int]RawDataRun, bytesPerCluster int64) (resolvedDataRuns map[int]DataRun) {

	return
}

func getDataRuns(dataRunBytes []byte, bytesPerCluster int64) (dataRuns map[int]DataRun) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	/*
		This function will parse the data runs from an MFT record.
		See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
	*/
	// Initialize a few variables
	rawDataRun := RawDataRun{}
	rawDataRuns := make(map[int]RawDataRun)
	offsetTracker := 0
	runCounter := 0

	for {
		if dataRunBytes[offsetTracker] == 0x00 {
			// Checks to see if we reached the end of the data runs. If so, break out of the loop.
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many bytes account for the
			// data run's offset and how many account for the data run's length.
			var offsetByteCount, lengthByteCount int
			byteToBeSplit := dataRunBytes[offsetTracker]
			offsetByteCount, lengthByteCount = getDataRunSplit(byteToBeSplit)
			if offsetByteCount == 0 && lengthByteCount == 0 {
				dataRuns = nil
				return
			}
			offsetTracker += 1

			// Pull out the the bytes that account for the data runs offset and length
			var lengthBytes, offsetBytes []byte

			lengthBytes = make([]byte, len(dataRunBytes[offsetTracker:(offsetTracker+lengthByteCount)]))
			copy(lengthBytes, dataRunBytes[offsetTracker:(offsetTracker+lengthByteCount)])
			offsetBytes = make([]byte, len(dataRunBytes[(offsetTracker+lengthByteCount):(offsetTracker+lengthByteCount+offsetByteCount)]))
			copy(offsetBytes, dataRunBytes[(offsetTracker+lengthByteCount):(offsetTracker+lengthByteCount+offsetByteCount)])

			// Convert the bytes for the data run offset and length to little endian int64
			rawDataRun.ClusterOffset = ConvertLittleEndianByteSliceToInt64(offsetBytes)
			if rawDataRun.ClusterOffset == 0 {
				dataRuns = nil
				return
			}

			rawDataRun.NumberOfClusters = ConvertLittleEndianByteSliceToInt64(lengthBytes)
			if rawDataRun.NumberOfClusters == 0 {
				dataRuns = nil
				return
			}
			// Append the data run to our data run struct
			rawDataRuns[runCounter] = rawDataRun

			// Increment the number order in preparation for the next data run.
			runCounter += 1

			// Set the offset tracker to the position of the next data run
			offsetTracker = offsetTracker + lengthByteCount + offsetByteCount
			if len(dataRunBytes) < offsetTracker {
				break
			}
		}
	}

	// Resolve Data Runs
	dataRuns = make(map[int]DataRun)
	offset := int64(0)
	for i := 0; i < len(rawDataRuns); i++ {
		offset = offset + (rawDataRuns[i].ClusterOffset * bytesPerCluster)
		dataRuns[i] = DataRun{
			AbsoluteOffset: offset,
			Length:         rawDataRuns[i].NumberOfClusters * bytesPerCluster,
		}
	}
	return
}

func getDataRunSplit(dataRunByte byte) (offsetByteCount, lengthByteCount int) {
	/*
		This function will split the first byte of a data run.
		See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
	*/
	// Convert the byte to a hex string
	hexToSplit := fmt.Sprintf("%x", dataRunByte)
	if len(hexToSplit) != 2 {
		offsetByteCount = 0
		lengthByteCount = 0
		return
	}

	// Split the hex string in half and return each half as an int
	offsetByteCount, err := strconv.Atoi(string(hexToSplit[0]))
	if err != nil {
		offsetByteCount = 0
		lengthByteCount = 0
		return
	}
	lengthByteCount, err = strconv.Atoi(string(hexToSplit[1]))
	if err != nil {
		offsetByteCount = 0
		lengthByteCount = 0
		return
	}
	return
}
