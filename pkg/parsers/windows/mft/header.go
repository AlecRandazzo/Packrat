// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

// RecordHeader contains parsed record header values.
type RecordHeader struct {
	AttributesOffset uint8
	RecordNumber     uint32
	Flags            Flags
}

// Flags contains parsed record header flag values.
type Flags struct {
	Deleted   bool
	Directory bool
}

var (
	headerAttributeOffsetLocation = byteshelper.NewDataLocation(0x14, 0x01)
	headerFlagLocation            = byteshelper.NewDataLocation(0x16, 0x01)
	headerRecordNumberLocation    = byteshelper.NewDataLocation(0x2C, 0x04)
)

// GetRecordHeaders parses the raw record header.
func GetRecordHeaders(input []byte) (RecordHeader, error) {
	// sanity checks
	size := len(input)
	if size == 0 {
		return RecordHeader{}, errors.New("received nil input")
	} else if size < 0x38 {
		return RecordHeader{}, fmt.Errorf("expected 38 input, instead it received %d", size)
	}

	// init return variable
	var header RecordHeader

	// get data
	buffer, _ := byteshelper.GetValue(input, headerAttributeOffsetLocation)
	header.AttributesOffset = buffer[0]

	buffer, _ = byteshelper.GetValue(input, headerFlagLocation)
	header.Flags = getHeaderFlags(buffer[0])

	buffer, _ = byteshelper.GetValue(input, headerRecordNumberLocation)
	header.RecordNumber = binary.LittleEndian.Uint32(buffer)

	return header, nil
}

const (
	codeDeletedFile = 0x00
	//codeActiveFile = 0x01
	//codeDeletedDirectory = 0x02
	codeDirectory = 0x03
)

// getHeaderFlags parses the raw record header flag.
func getHeaderFlags(input byte) Flags {
	// init return variable
	var flags Flags

	// interpret byte
	switch input {
	case codeDeletedFile:
		flags.Deleted = true
	case codeDirectory:
		flags.Directory = true
	}

	return flags
}

var (
	MagicNumber                  = []byte{0x46, 0x49, 0x4c, 0x45, 0x30} // FILE0
	mftRecordMagicNumberLocation = byteshelper.NewDataLocation(0x00, 0x05)
)

// ValidateMftRecordBytes quickly checks to see if the raw mft record is a valid. It will return no error if the record is value.
func ValidateMftRecordBytes(input []byte) error {
	// sanity checks
	size := len(input)
	if size == 0 {
		return errors.New("received nil input")
	}
	if size < 0x05 {
		return errors.New("received less than 4 input")
	}

	// check magic number
	buffer, _ := byteshelper.GetValue(input, mftRecordMagicNumberLocation)

	if bytes.Compare(buffer, MagicNumber) != 0 {
		return errors.New("incorrect magic number")
	}

	return nil
}
