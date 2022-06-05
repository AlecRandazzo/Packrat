// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

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

var (
	dataAttributeResidentFlagLocation = byteshelper.NewDataLocation(0x08, 0x01)
)

// getDataAttribute parses the raw data attribute. bytesPerCluster is typically 4096. The function returns an interface type that is either a ResidentDataAttribute or NonResidentDataAttribute.
func getDataAttribute(input []byte, bytesPerCluster uint) (dataAttribute interface{}, err error) {
	// Sanity checks on data the method receives to make sure it can successfully do work on the data.
	sizeOfRawDataAttribute := len(input)
	if sizeOfRawDataAttribute == 0 {
		return nil, errors.New("received nil input")
	} else if sizeOfRawDataAttribute <= 0x08 {
		return nil, errors.New("received input less than 8")
	}
	if bytesPerCluster == 0 {
		return nil, errors.New("did not receive a value for input per cluster")
	}

	var buffer []byte
	// Check to see if the attribute is resident or not. Parses the data accordingly.
	buffer, err = byteshelper.GetValue(input, dataAttributeResidentFlagLocation)
	if buffer[0] == 0x00 {
		var residentDataAttribute ResidentDataAttribute
		residentDataAttribute, err = getResidentDataAttribute(input)
		if err != nil {
			return nil, fmt.Errorf("failed to parse resident data attribute: %w", err)
		}
		return residentDataAttribute, nil
	}
	var nonResidentDataAttribute NonResidentDataAttribute
	nonResidentDataAttribute, err = getNonResidentDataAttribute(input, bytesPerCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to parse non resident data attribute: %w", err)
	}

	return nonResidentDataAttribute, nil
}

const offsetResidentData = 0x18

// getResidentDataAttribute parses the raw resident data attribute receiver and returns the resident data attribute input.
func getResidentDataAttribute(input []byte) (ResidentDataAttribute, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return ResidentDataAttribute{}, errors.New("received nil input")
	} else if size < offsetResidentData {
		return ResidentDataAttribute{}, fmt.Errorf("expected to receive at least 18 input, but received %d", size)
	}

	// Copy resident data
	residentDataAttribute := make(ResidentDataAttribute, len(input[offsetResidentData:]))
	copy(residentDataAttribute, input[offsetResidentData:])
	return residentDataAttribute, nil
}

const offsetDataRunOffset = 0x20

// getNonResidentDataAttribute parses the raw non resident data attribute.
func getNonResidentDataAttribute(input []byte, bytesPerCluster uint) (NonResidentDataAttribute, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return NonResidentDataAttribute{}, errors.New("received nil input")
	} else if size <= offsetDataRunOffset {
		return NonResidentDataAttribute{}, fmt.Errorf("expected to receive at least 18 input, but received %d", size)
	}

	// Identify offset of the data runs in the data Attribute
	dataRunOffset := input[offsetDataRunOffset]

	// Verify we aren't going outside the bounds of the byte slice
	if size < int(dataRunOffset) {
		return NonResidentDataAttribute{}, errors.New("data run offset is beyond the size of the byte slice")
	}

	// Send the input to be parsed
	var dataAttribute NonResidentDataAttribute
	dataAttribute.DataRuns, _ = getDataRuns(input[dataRunOffset:], bytesPerCluster)

	return dataAttribute, nil
}

// Parse parses the raw data run receiver and returns data runs. The input per cluster argument is used to calculate data run information.
func getDataRuns(input []byte, bytesPerCluster uint) (DataRuns, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return DataRuns{}, errors.New("received nil input")
	}

	// Initialize a few variables
	UnresolvedDataRuns := make(unresolvedDataRuns)
	dataRuns := make(DataRuns, 0)
	offset := 0
	runCounter := 0

	for {
		// Checks to see if we reached the end of the data runs. If so, break out of the loop.
		if size <= offset {
			break
		} else if input[offset] == 0x00 {
			break
		} else {
			// Take the first byte of a data run and send it to get split so we know how many input account for the
			// data run's offset and how many account for the data run's length.
			dataRunSplit := splitDataRunByte(input[offset])
			offset++

			// Pull out the input that account for the data runs offset and length
			var lengthBytes, offsetBytes []byte
			if offset+dataRunSplit.lengthByteCount > size {
				break
			}
			lengthBytes = make([]byte, len(input[offset:(offset+dataRunSplit.lengthByteCount)]))
			copy(lengthBytes, input[offset:(offset+dataRunSplit.lengthByteCount)])
			if offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount > size {
				break
			}
			offsetBytes = make([]byte, len(input[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)]))
			copy(offsetBytes, input[(offset+dataRunSplit.lengthByteCount):(offset+dataRunSplit.lengthByteCount+dataRunSplit.offsetByteCount)])

			// Convert the input for the data run offset and length to little endian int64
			var UnresolvedDataRun unresolvedDataRun
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
		dataRunOffset = dataRunOffset + (UnresolvedDataRuns[i].clusterOffset * int64(bytesPerCluster))
		dataRuns[i] = DataRun{
			AbsoluteOffset: dataRunOffset,
			Length:         UnresolvedDataRuns[i].numberOfClusters * int64(bytesPerCluster),
		}
	}

	return dataRuns, nil
}

// splitDataRunByte will split the first byte of a data run.
// See the following for a good write up on data runs: https://homepage.cs.uri.edu/~thenry/csc487/video/66_NTFS_Data_Runs.pdf
func splitDataRunByte(input byte) dataRunSplit {
	// init return variable
	var split dataRunSplit

	// Convert the byte to a hex string
	hexToSplit := fmt.Sprintf("%x", input)

	if len(hexToSplit) == 1 {
		split.offsetByteCount = 0
		split.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
	} else {
		// Split the hex string in half and return each half as an int
		split.offsetByteCount, _ = strconv.Atoi(string(hexToSplit[0]))
		split.lengthByteCount, _ = strconv.Atoi(string(hexToSplit[1]))
	}

	return split
}
