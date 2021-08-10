// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

// RawDataAttribute is an alias for a raw data attribute. Used as a receiver to the parse() method.
type RawDataAttribute []byte

// RawResidentDataAttribute is an alias for a raw resident data attribute. Used as a receiver to the parse() method. We don't really do anything this this currently.
type RawResidentDataAttribute []byte

// RawNonResidentDataAttribute is an alias for a raw nonresident data attribute. Used as a receiver to the parse() method.
type RawNonResidentDataAttribute []byte

// RawDataRuns is an alias for a raw data runs. Used as a receiver to the parse() method.
// See this for more details on data runs: https://flatcap.org/linux-ntfs/ntfs/concepts/data_runs.html
type RawDataRuns []byte
type rawDataRunSplitByte byte

// ResidentDataAttribute is an alias for a resident data attribute.
type ResidentDataAttribute []byte

// NonResidentDataAttribute is an alias for a parsed non-resident data attribute.
type NonResidentDataAttribute struct {
	DataRuns DataRuns
}

type unresolvedDataRun struct {
	numberOrder      int
	clusterOffset    int64
	numberOfClusters int64
}

type unresolvedDataRuns map[int]unresolvedDataRun

// DataRuns contains an ordered slice of parsed data runs
type DataRuns map[int]DataRun

// DataRun contains a parsed data run which contains the absolute offset of where the data run resides in the volume and the length of the data run.
type DataRun struct {
	AbsoluteOffset int64
	Length         int64
}

type dataRunSplit struct {
	offsetByteCount int
	lengthByteCount int
}

// DataAttribute contains information about a parsed data attribute.
type DataAttribute struct {
	TotalSize                uint8
	FlagResident             bool
	ResidentDataAttribute    ResidentDataAttribute
	NonResidentDataAttribute NonResidentDataAttribute
}

// Parse parses the raw data attribute receiver and returns a non resident data attribute or a resident data attribute. The bytes per cluster argument is used to calculate data run information.
func (rawDataAttribute RawDataAttribute) Parse(bytesPerCluster int64) (nonResidentDataAttribute NonResidentDataAttribute, residentDataAttribute ResidentDataAttribute, err error) {

	// Sanity checks on data the method receives to make sure it can successfully do work on the data.
	const offsetResidentFlag = 0x08
	sizeOfRawDataAttribute := len(rawDataAttribute)
	if sizeOfRawDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawDataAttribute <= offsetResidentFlag {
		err = errors.New("received bytes less than 8")
		return
	}
	if bytesPerCluster == 0 {
		err = errors.New("did not receive a value for bytes per cluster")
		return
	}

	// Check to see if the attribute is resident or not. Parses the data accordingly.
	if rawDataAttribute[offsetResidentFlag] == 0x00 {
		rawResidentDataAttribute := RawResidentDataAttribute(make([]byte, sizeOfRawDataAttribute))
		copy(rawResidentDataAttribute, rawDataAttribute)
		residentDataAttribute, err = rawResidentDataAttribute.Parse()
		if err != nil {
			err = fmt.Errorf("failed to parse resident data attribute: %w", err)
			return
		}
		return
	}

	rawNonResidentDataAttribute := RawNonResidentDataAttribute(make([]byte, sizeOfRawDataAttribute))
	copy(rawNonResidentDataAttribute, rawDataAttribute)
	nonResidentDataAttribute, err = rawNonResidentDataAttribute.Parse(bytesPerCluster)
	if err != nil {
		err = fmt.Errorf("failed to parse non resident data attribute: %w", err)
		return
	}

	return
}

// Parse parses the raw resident data attribute receiver and returns the resident data attribute bytes.
func (rawResidentDataAttribute RawResidentDataAttribute) Parse() (residentDataAttribute ResidentDataAttribute, err error) {
	// Sanity check to make sure the method received good data
	const offsetResidentData = 0x18
	sizeOfRawResidentDataAttribute := len(rawResidentDataAttribute)
	if sizeOfRawResidentDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawResidentDataAttribute < offsetResidentData {
		err = fmt.Errorf("expected to receive at least 18 bytes, but received %d", sizeOfRawResidentDataAttribute)
		return
	}
	sizeOfResidentDataAttribute := len(rawResidentDataAttribute[offsetResidentData:])
	residentDataAttribute = make(ResidentDataAttribute, sizeOfResidentDataAttribute)
	copy(residentDataAttribute, rawResidentDataAttribute[offsetResidentData:])
	return
}

// Parse parses the raw non resident data attribute receiver and returns a non resident data attribute. The bytes per cluster argument is used to calculate data run information.
func (rawNonResidentDataAttribute RawNonResidentDataAttribute) Parse(bytesPerCluster int64) (nonResidentDataAttributes NonResidentDataAttribute, err error) {
	// Sanity check to make sure the method received good data
	const offsetDataRunOffset = 0x20
	sizeOfRawNonResidentDataAttribute := len(rawNonResidentDataAttribute)
	if sizeOfRawNonResidentDataAttribute == 0 {
		err = errors.New("received nil bytes")
		return
	} else if sizeOfRawNonResidentDataAttribute <= offsetDataRunOffset {
		err = fmt.Errorf("expected to receive at least 18 bytes, but received %d", sizeOfRawNonResidentDataAttribute)
		return
	}

	// Identify offset of the data runs in the data Attribute
	dataRunOffset := rawNonResidentDataAttribute[offsetDataRunOffset]

	// Verify we aren't going outside the bounds of the byte slice
	if sizeOfRawNonResidentDataAttribute < int(dataRunOffset) {
		err = errors.New("data run offset is beyond the size of the byte slice")
		return
	}

	// Pull out the data run bytes
	rawDataRuns := RawDataRuns(make([]byte, sizeOfRawNonResidentDataAttribute))
	copy(rawDataRuns, rawNonResidentDataAttribute[dataRunOffset:])

	// Send the bytes to be parsed
	nonResidentDataAttributes.DataRuns, _ = rawDataRuns.Parse(bytesPerCluster)

	return
}

// Parse parses the raw data run receiver and returns data runs. The bytes per cluster argument is used to calculate data run information.
func (rawDataRuns RawDataRuns) Parse(bytesPerCluster int64) (dataRuns DataRuns, err error) {
	// Sanity check that the method received good data
	if rawDataRuns == nil {
		err = errors.New("received null bytes")
		return
	}

	// Initialize a few variables
	UnresolvedDataRun := unresolvedDataRun{}
	UnresolvedDataRuns := make(unresolvedDataRuns)
	sizeOfRawDataRuns := len(rawDataRuns)
	dataRuns = make(DataRuns)
	offset := 0
	runCounter := 0

	for {
		// Checks to see if we reached the end of the data runs. If so, break out of the loop.
		if rawDataRuns[offset] == 0x00 || sizeOfRawDataRuns < offset {
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many bytes account for the
			// data run's offset and how many account for the data run's length.
			byteToBeSplit := rawDataRunSplitByte(rawDataRuns[offset])
			dataRunSplit := byteToBeSplit.parse()
			offset++

			// Pull out the the bytes that account for the data runs offset2 and length
			var lengthBytes, offsetBytes []byte

			lengthBytes = make([]byte, len(rawDataRuns[offset:(offset+dataRunSplit.lengthByteCount)]))
			copy(lengthBytes, rawDataRuns[offset:(offset+dataRunSplit.lengthByteCount)])
			offsetBytes = make([]byte, len(rawDataRuns[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)]))
			copy(offsetBytes, rawDataRuns[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)])

			// Convert the bytes for the data run offset and length to little endian int64
			UnresolvedDataRun.clusterOffset, _ = byteshelper.LittleEndianBinaryToInt64(offsetBytes)
			UnresolvedDataRun.numberOfClusters, _ = byteshelper.LittleEndianBinaryToInt64(lengthBytes)

			// Append the data run to our data run struct
			UnresolvedDataRuns[runCounter] = UnresolvedDataRun

			// Increment the number order in preparation for the next data run.
			runCounter++

			// Set the offset tracker to the position of the next data run
			offset = offset + dataRunSplit.lengthByteCount + dataRunSplit.offsetByteCount
		}
	}

	// Resolve Data Runs
	dataRunOffset := int64(0)
	for i := 0; i < len(UnresolvedDataRuns); i++ {
		dataRunOffset = dataRunOffset + (UnresolvedDataRuns[i].clusterOffset * bytesPerCluster)
		dataRuns[i] = DataRun{
			AbsoluteOffset: dataRunOffset,
			Length:         UnresolvedDataRuns[i].numberOfClusters * bytesPerCluster,
		}
	}
	return
}

// This function will split the first byte of a data run.
// See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
func (rawDataRunSplitByte rawDataRunSplitByte) parse() (dataRunSplit dataRunSplit) {
	// Convert the byte to a hex string
	hexToSplit := fmt.Sprintf("%x", rawDataRunSplitByte)

	if len(hexToSplit) == 1 {
		dataRunSplit.offsetByteCount = 0
		dataRunSplit.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
	} else {
		// Split the hex string in half and return each half as an int
		dataRunSplit.offsetByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
		dataRunSplit.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[1]))
	}
	return
}
