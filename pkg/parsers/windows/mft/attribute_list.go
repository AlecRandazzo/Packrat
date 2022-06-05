// Copyright (c) 2022 Alec Randazzo

package mft

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/AlecRandazzo/Packrat/pkg/parsers/general/byteshelper"
)

// AttributeListAttribute contains information about a attribute list attribute
type AttributeListAttribute struct {
	Type                     byte
	MFTReferenceRecordNumber uint32
}

// AttributeListAttributes is an array of AttributeListAttribute
type AttributeListAttributes []AttributeListAttribute

var (
	attributeListTypeLocation         = byteshelper.NewDataLocation(0x00, 0x01)
	attributeListRecordLengthLocation = byteshelper.NewDataLocation(offsetRecordLength, lengthRecordLength)
)

const (
	offsetRecordLength = 0x04
	lengthRecordLength = 0x02

	offsetFirstSubAttribute = 0x18

	offsetMFTReferenceRecordNumber = 0x10
	lengthMFTReferenceRecordNumber = 0x04
)

// getAttributeListAttribute parses a raw attribute list attribute
func getAttributeListAttribute(input []byte) (AttributeListAttributes, error) {
	// Sanity checking
	size := len(input)
	if size == 0 {
		return AttributeListAttributes{}, errors.New("received nil input")
	}

	buffer, _ := byteshelper.GetValue(input, attributeListTypeLocation)
	if buffer[0] != 0x20 {
		return AttributeListAttributes{}, fmt.Errorf("receive an attribute thats not an attribute list. Attribute magic number is %x", buffer[0])
	}

	buffer, _ = byteshelper.GetValue(input, attributeListRecordLengthLocation)
	recordLength := binary.LittleEndian.Uint16(buffer)
	if int(recordLength) != size {
		return AttributeListAttributes{}, fmt.Errorf("received a byte slice thats not equal to the expected attribute length. Size received was %d but expected %d", size, recordLength)
	}

	attributeList := make(AttributeListAttributes, 0)

	pointerToSubAttribute := offsetFirstSubAttribute
	for pointerToSubAttribute < size {
		err := validateAttribute(input[pointerToSubAttribute])
		if err != nil {
			return attributeList, nil
		}

		sizeOfSubAttribute := binary.LittleEndian.Uint16(input[pointerToSubAttribute+offsetRecordLength : pointerToSubAttribute+offsetRecordLength+lengthRecordLength])

		recordNumber := binary.LittleEndian.Uint32(input[pointerToSubAttribute+offsetMFTReferenceRecordNumber : pointerToSubAttribute+offsetMFTReferenceRecordNumber+lengthMFTReferenceRecordNumber])
		attributeListAttribute := AttributeListAttribute{
			Type:                     input[pointerToSubAttribute],
			MFTReferenceRecordNumber: recordNumber,
		}
		attributeList = append(attributeList, attributeListAttribute)
		pointerToSubAttribute += int(sizeOfSubAttribute)
	}

	return attributeList, nil
}
