// Copyright (c) 2020 Alec Randazzo

package mft

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// MasterFileTableRecord contains information on a parsed MFT record
type MasterFileTableRecord struct {
	RecordHeader                  RecordHeader
	StandardInformationAttributes StandardInformationAttribute
	FileNameAttributes            []FileNameAttribute
	DataAttribute                 DataAttribute
	AttributeList                 AttributeListAttributes
}

type write interface {
	Write(p []byte) (n int, err error)
}

func Parse(reader io.Reader, writer io.Writer) {

}

// ParseMftRecord the raw MFT record receiver and returns a parsed mft record.
func ParseMftRecord(input []byte, bytesPerCluster int64) (MasterFileTableRecord, error) {
	// Sanity checks
	size := len(input)
	if size == 0 {
		return MasterFileTableRecord{}, errors.New("received nil input")
	}
	if bytesPerCluster == 0 {
		return MasterFileTableRecord{}, errors.New("input per cluster of 0, typically this value is 4096")
	}

	// init return variables
	var mft MasterFileTableRecord

	err := ValidateMftRecordBytes(input)
	if err != nil {
		return MasterFileTableRecord{}, fmt.Errorf("this is not an mft record: %w", err)
	}

	input = trimSlackSpace(input)

	mft.RecordHeader, err = GetRecordHeaders(input)
	if err != nil {
		return MasterFileTableRecord{}, fmt.Errorf("failed to get record headers: %w", err)
	}

	if mft.RecordHeader.RecordNumber == 419091 {
		fmt.Printf("%x\n", input)
	}

	var rawAttributes [][]byte
	rawAttributes, err = GetRawAttributes(input, mft.RecordHeader)
	if err != nil {
		return MasterFileTableRecord{}, fmt.Errorf("failed to get raw data attributes: %w", err)
	}

	mft.FileNameAttributes, mft.StandardInformationAttributes, mft.DataAttribute, mft.AttributeList, _ = GetAttributes(rawAttributes, bytesPerCluster)
	return mft, nil
}

// Trims off slack space after end sequence 0xffffffff
func trimSlackSpace(input []byte) []byte {
	lenMftRecordBytes := len(input)
	mftRecordEndByteSequence := []byte{0xff, 0xff, 0xff, 0xff}
	for i := 0; i < (lenMftRecordBytes - 4); i++ {
		if bytes.Equal(input[i:i+0x04], mftRecordEndByteSequence) {
			input = input[:i]
			break
		}
	}

	return input
}
